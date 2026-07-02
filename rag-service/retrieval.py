import os
import json
from openai import OpenAI
import qdrant_client
from llama_index.core import VectorStoreIndex
from llama_index.vector_stores.qdrant import QdrantVectorStore
from llama_index.embeddings.fastembed import FastEmbedEmbedding
from ingestion import get_qdrant_client, get_collection_name
import logging

logger = logging.getLogger(__name__)

def retrieve_and_generate(query: str, doc_id: str, user_id: str, chat_history: list):
    qdrant_url = os.environ.get("QDRANT_URL", "http://localhost:6333")
    client = get_qdrant_client()
    collection_name = get_collection_name(user_id, doc_id)
    
    # 1. Verify Qdrant collection exists
    try:
        client.get_collection(collection_name)
    except Exception as e:
        logger.error(f"Qdrant collection {collection_name} does not exist: {str(e)}")
        yield f"event: error\ndata: {json.dumps({'message': 'Document index not found. Please wait for processing to finish.'})}\n\n"
        return

    # 2. Setup retriever
    vector_store = QdrantVectorStore(client=client, collection_name=collection_name)
    embed_model = FastEmbedEmbedding(model_name="BAAI/bge-small-en-v1.5")
    index = VectorStoreIndex.from_vector_store(vector_store, embed_model=embed_model)
    
    # Retrieve top 5 chunks
    retriever = index.as_retriever(similarity_top_k=5)
    
    try:
        retrieved_nodes = retriever.retrieve(query)
    except Exception as e:
        logger.error(f"Error retrieving nodes from Qdrant: {str(e)}")
        yield f"event: error\ndata: {json.dumps({'message': 'Failed to retrieve information from document.'})}\n\n"
        return

    # 3. Format sources
    sources = []
    for node_with_score in retrieved_nodes:
        node = node_with_score.node
        sources.append({
            "text": node.text,
            "page_number": node.metadata.get("page_number", 0)
        })
        
    logger.info(f"Retrieved {len(sources)} source chunks for query")

    # 4. Stream sources metadata to the client first
    yield f"event: sources\ndata: {json.dumps(sources)}\n\n"

    # 5. Build prompts
    context_str = "\n\n".join([f"[Page {src['page_number']}]: {src['text']}" for src in sources])
    
    system_prompt = (
        "You are FitMind, a professional and compassionate AI-powered personal health advisor. "
        "Your goal is to answer the user's health questions accurately based ONLY on the provided health documents "
        "(blood reports, diet plans, workout logs, medical summaries) and chat history.\n\n"
        "Guidelines:\n"
        "1. Rely only on the provided Context to answer. If the context doesn't contain the answer, say clearly "
        "that you don't have enough information from the documents to answer, but still provide helpful general advice "
        "based on the query.\n"
        "2. Reference the page numbers in your response if relevant, but do not make up any information.\n"
        "3. Maintain a warm, encouraging, and medical-advisor tone.\n\n"
        "CONTEXT FROM USER'S DOCUMENTS:\n"
        f"{context_str}"
    )

    messages = [{"role": "system", "content": system_prompt}]
    for msg in chat_history:
        messages.append({"role": msg.role, "content": msg.content})
        
    # Append latest query
    messages.append({"role": "user", "content": query})

    # 6. Stream OpenAI Chat Completion
    api_key = os.environ.get("OPENAI_API_KEY")
    if not api_key:
        yield f"event: error\ndata: {json.dumps({'message': 'OpenAI API key is missing in configuration.'})}\n\n"
        return
        
    try:
        openai_client = OpenAI(
            api_key=api_key,
            base_url="https://api.groq.com/openai/v1"
        )
        response = openai_client.chat.completions.create(
            model="llama3-8b-8192",
            messages=messages,
            stream=True
        )
        
        for chunk in response:
            token = chunk.choices[0].delta.content
            if token:
                yield f"event: text\ndata: {json.dumps(token)}\n\n"
                
        yield "event: done\ndata: {}\n\n"
        
    except Exception as e:
        logger.error(f"Error streaming response from OpenAI: {str(e)}")
        yield f"event: error\ndata: {json.dumps({'message': f'Failed to generate response: {str(e)}'})}\n\n"
