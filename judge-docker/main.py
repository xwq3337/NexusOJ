from pydantic import BaseModel
from fastapi import FastAPI
import uvicorn
from threading import Thread
from Judger.cpp.judge import C_cpp
from Judger.python.judge import Python
from Judger.rust.judge import Rust
from Judger.javascript.judge import JavaScript
from Judger.go.judge import Go
app = FastAPI()
def run_app(app, host: str, port: int):
    uvicorn.run(app, host=host, port=port)

class JudgeModel(BaseModel):
    code : str
    language : str
    input : str
    timeLimit: int
    memoryLimit: int
    args: str
@app.get('/language')
async def getLanguage():
    return ["c","c++","python","javascript","go","rust"]

@app.post('/judge')
async def judge(data :JudgeModel):
    if data.language == "cpp" or data.language == "c":
        return await C_cpp(data.code, data.language, data.input)
    elif data.language == "python":
        return await Python(data.code,data.language,data.input)
    elif data.language == "rust":
        return await Rust(data.code,data.language,data.input)
    elif data.language == "javascript":
        return await JavaScript(data.code,data.language,data.input)
    elif data.language == "go":
        return await Go(data.code,data.language,data.input)

IP_ADDR = "0.0.0.0"
if __name__ == "__main__":
    Thread(target=run_app, args=(app, IP_ADDR, 28000)).start()
    Thread(target=run_app, args=(app, IP_ADDR, 28001)).start()
    Thread(target=run_app, args=(app, IP_ADDR, 28002)).start()
    Thread(target=run_app, args=(app, IP_ADDR, 28003)).start()
