import time
import os
import signal
import subprocess
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler
 
class ReloadHandler(FileSystemEventHandler):
    def __init__(self):
        self.proc = None 
        self.reload_process()
    def on_modified(self, event):
        if event.src_path.endswith('.py'):
            print(f"Detected change in {event.src_path}, reloading...")
            self.reload_process()

    def reload_process(self):
        if self.proc is not None:
            self.proc.send_signal(signal.SIGTERM)
            self.proc.wait()
            
        self.proc = subprocess.Popen(["python3", "main.py"])

if __name__ == "__main__":
    path = '.'
    event_handler = ReloadHandler()
    observer = Observer()
    observer.schedule(event_handler, path, recursive=True)
    observer.start()
    
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        observer.stop()
        if event_handler.proc is not None:
            event_handler.proc.terminate()  # 确保在退出时终止子进程
            event_handler.proc.wait()
    observer.join()