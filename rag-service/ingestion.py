import os
import fitz  # PyMuPDF
from supabase import create_client, Client
import qdrant_client
from qdrant_client.models import Distance, VectorParams
from llama_index.core import Document, StorageContext, VectorStoreIndex
from llama_index.core.node_parser import TokenTextSplitter
from llama_index.embeddings.openai import OpenAIEmbedding
from llama_index.vector_stores.qdrant import QdrantVectorStore
import logging

logger = logging.getLogger(__name__)

def get_supabase_client() -> Client:
    supabase_url = os.environ.get("SUPABASE_URL")
    supabase_key = os.environ.get("SUPABASE_SERVICE_ROLE_KEY")
    if not supabase_url or not supabase_key:
        raise ValueError("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set in environment variables")
    return create_client(supabase_url, supabase_key)

def get_qdrant_client() -> qdrant_client.QdrantClient:
    qdrant_url = os.environ.get("QDRANT_URL", "http://localhost:6333")
    return qdrant_client.QdrantClient(url=qdrant_url)

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
            vectors_config=VectorParams(size=1536, distance=Distance.COSINE),
        )
    
    # 5. Build and Store Index
    vector_store = QdrantVectorStore(client=client, collection_name=collection_name)
    storage_context = StorageContext.from_defaults(vector_store=vector_store)
    embed_model = OpenAIEmbedding(model="text-embedding-3-small")
    
    logger.info(f"Indexing chunks in Qdrant collection {collection_name}")
    index = VectorStoreIndex(
        nodes,
        storage_context=storage_context,
        embed_model=embed_model
    )
    logger.info(f"Successfully ingested and indexed doc_id={doc_id} for user_id={user_id}")
    return len(nodes)
