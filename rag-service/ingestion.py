import os
import fitz  # PyMuPDF
from supabase import create_client, Client
import qdrant_client
from qdrant_client.models import Distance, VectorParams
from llama_index.core import Document, StorageContext, VectorStoreIndex
from llama_index.core.node_parser import TokenTextSplitter
from llama_index.embeddings.fastembed import FastEmbedEmbedding
from llama_index.vector_stores.qdrant import QdrantVectorStore
import logging
from analyzer import analyze_document_text


logger = logging.getLogger(__name__)

def get_supabase_client() -> Client:
    supabase_url = os.environ.get("SUPABASE_URL")
    supabase_key = os.environ.get("SUPABASE_SERVICE_ROLE_KEY")
    if not supabase_url or not supabase_key:
        raise ValueError("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set in environment variables")
    return create_client(supabase_url, supabase_key)

def get_qdrant_client() -> qdrant_client.QdrantClient:
    qdrant_url = os.environ.get("QDRANT_URL")
    url = qdrant_url or "http://localhost:6333"
    
    # Try connecting to the server. If it's unreachable, fall back to file-backed Qdrant.
    try:
        client = qdrant_client.QdrantClient(url=url, timeout=2.0)
        client.get_collections() # test connection
        logger.info(f"Successfully connected to Qdrant server at {url}")
        return client
    except Exception:
        db_path = os.path.join(os.path.dirname(os.path.abspath(__file__)), "qdrant_db")
        logger.info(f"Qdrant server at {url} is unreachable. Falling back to local file-backed Qdrant DB at: {db_path}")
        return qdrant_client.QdrantClient(path=db_path)


def get_collection_name(user_id: str, doc_id: str) -> str:
    # Ensure name contains only alphanumeric, hyphens or underscores
    clean_user = user_id.replace("-", "_")
    clean_doc = doc_id.replace("-", "_")
    return f"col_{clean_user}_{clean_doc}"

def download_file_from_supabase(storage_path: str) -> bytes:
    supabase = get_supabase_client()
    # Storage bucket name defaults to "documents" as per the plan
    # The storage path is e.g. "uploads/user_id/doc_id.pdf"
    bucket_name = "documents"
    logger.info(f"Downloading storage path: {storage_path} from bucket: {bucket_name}")
    try:
        file_data = supabase.storage.from_(bucket_name).download(storage_path)
        return file_data
    except Exception as e:
        logger.error(f"Error downloading file from Supabase: {str(e)}")
        raise e

def parse_pdf(file_bytes: bytes, user_id: str, doc_id: str) -> list[Document]:
    logger.info("Parsing PDF with PyMuPDF")
    doc = fitz.open(stream=file_bytes, filetype="pdf")
    documents_list = []
    
    for page_idx, page in enumerate(doc):
        text = page.get_text()
        page_num = page_idx + 1
        documents_list.append(Document(
            text=text,
            metadata={
                "page_number": page_num,
                "doc_id": doc_id,
                "user_id": user_id
            }
        ))
    logger.info(f"Parsed {len(documents_list)} pages from PDF")
    return documents_list

def ingest_document(user_id: str, doc_id: str, storage_path: str):
    logger.info(f"Starting ingestion for user_id={user_id}, doc_id={doc_id}")
    
    # 1. Download file
    file_bytes = download_file_from_supabase(storage_path)
    
    # 2. Parse PDF
    documents_list = parse_pdf(file_bytes, user_id, doc_id)
    
    # 3. Split text
    # Chunk size 512, overlap 50 as specified
    splitter = TokenTextSplitter(chunk_size=512, chunk_overlap=50)
    nodes = splitter.get_nodes_from_documents(documents_list)
    logger.info(f"Split documents into {len(nodes)} chunks")
    
    # 4. Setup Qdrant Collection
    client = get_qdrant_client()
    collection_name = get_collection_name(user_id, doc_id)
    
    try:
        client.get_collection(collection_name)
        logger.info(f"Collection {collection_name} already exists. It will be reused.")
    except Exception:
        logger.info(f"Collection {collection_name} does not exist. Creating it.")
        client.create_collection(
            collection_name=collection_name,
            vectors_config=VectorParams(size=384, distance=Distance.COSINE),
        )
    
    # 5. Build and Store Index
    vector_store = QdrantVectorStore(client=client, collection_name=collection_name)
    storage_context = StorageContext.from_defaults(vector_store=vector_store)
    embed_model = FastEmbedEmbedding(model_name="BAAI/bge-small-en-v1.5")
    
    logger.info(f"Indexing chunks in Qdrant collection {collection_name}")
    index = VectorStoreIndex(
        nodes,
        storage_context=storage_context,
        embed_model=embed_model
    )
    logger.info(f"Successfully ingested and indexed doc_id={doc_id} for user_id={user_id}")
    
    # 6. Automatically run analysis
    try:
        logger.info(f"Running automated analysis for doc_id={doc_id}")
        full_text = "\n".join([doc.text for doc in documents_list])
        analysis_result = analyze_document_text(full_text)
        
        supabase = get_supabase_client()
        
        # Save analysis result
        analysis_data = {
            "doc_id": doc_id,
            "user_id": user_id,
            "metrics": analysis_result.get("metrics", []),
            "summary": analysis_result.get("summary", ""),
            "overall_status": analysis_result.get("overall_status", "good")
        }
        
        logger.info(f"Saving analysis result to Supabase for doc_id={doc_id}")
        supabase.table("document_analyses").insert(analysis_data).execute()
        
        # Update document status to 'analyzed'
        logger.info(f"Updating document status to 'analyzed' for doc_id={doc_id}")
        supabase.table("documents").update({"status": "analyzed"}).eq("id", doc_id).execute()
        
    except Exception as e:
        logger.error(f"Failed to automatically analyze document: {str(e)}")
        # Even if analysis fails, vector indexing succeeded. We set it to ready.
        try:
            supabase = get_supabase_client()
            supabase.table("documents").update({"status": "ready"}).eq("id", doc_id).execute()
        except Exception as db_err:
            logger.error(f"Failed to update document status to ready: {str(db_err)}")
            
    return len(nodes)

