#!/usr/bin/env python3
"""
Local preview server for Pomme landing page
"""

import http.server
import socketserver
import os
import sys
import webbrowser
from threading import Timer
import subprocess
import signal
import time

PORT = 8000
DIRECTORY = os.path.dirname(os.path.abspath(__file__))

class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)
    
    def end_headers(self):
        # Add headers to prevent caching during development
        self.send_header('Cache-Control', 'no-store, no-cache, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')
        super().end_headers()
    
    def log_message(self, format, *args):
        # Custom log format
        sys.stderr.write(f"[{self.log_date_time_string()}] {format % args}\n")

def open_browser():
    webbrowser.open(f'http://localhost:{PORT}')

def kill_port_process(port):
    """Kill any process using the specified port"""
    killed = False
    try:
        # Find process using the port
        result = subprocess.run(['lsof', '-t', f'-i:{port}'], 
                              capture_output=True, text=True)
        if result.stdout.strip():
            pids = result.stdout.strip().split('\n')
            for pid in pids:
                try:
                    os.kill(int(pid), signal.SIGTERM)
                    print(f"✓ Killed process {pid} using port {port}")
                    killed = True
                except:
                    pass
    except:
        # lsof might not be available, try fuser as fallback
        try:
            result = subprocess.run(['fuser', '-k', f'{port}/tcp'], 
                         capture_output=True, stderr=subprocess.DEVNULL)
            if result.returncode == 0:
                killed = True
        except:
            pass
    
    # Give the OS time to release the port
    if killed:
        time.sleep(1)

def main():
    os.chdir(DIRECTORY)
    
    # Kill any existing process on the port
    kill_port_process(PORT)
    
    # Try to bind to the port with retries
    max_retries = 3
    retry_count = 0
    
    while retry_count < max_retries:
        try:
            with socketserver.TCPServer(("", PORT), Handler) as httpd:
                print(f"\n🍎 Pomme Landing Page Preview Server")
                print(f"{'─' * 40}")
                print(f"🌐 Serving at: http://localhost:{PORT}")
                print(f"📁 Directory: {DIRECTORY}")
                print(f"{'─' * 40}")
                print(f"Press Ctrl+C to stop the server\n")
                
                # Open browser after a short delay
                Timer(0.5, open_browser).start()
                
                try:
                    httpd.serve_forever()
                except KeyboardInterrupt:
                    print("\n\n✨ Server stopped")
                    sys.exit(0)
                break
        except OSError as e:
            if e.errno == 98:  # Address already in use
                retry_count += 1
                if retry_count < max_retries:
                    print(f"⚠️  Port {PORT} still in use, retrying in 2 seconds...")
                    time.sleep(2)
                    kill_port_process(PORT)
                else:
                    print(f"❌ Failed to start server: Port {PORT} is still in use after {max_retries} attempts")
                    print("   Try killing the process manually or use a different port")
                    sys.exit(1)
            else:
                raise

if __name__ == "__main__":
    main()