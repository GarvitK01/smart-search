import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer

app = FastAPI()
model = SentenceTransformer("./model")


class EmbedResponse(BaseModel):
    embeddings: list[list[float]]

class EmbedRequest(BaseModel):
    texts: list[str]

@app.post("/embed")
def embed(request: EmbedRequest):
    embeddings = model.encode(request.texts)
    return EmbedResponse(embeddings=embeddings.tolist())

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)