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

def main():
    os.chdir(DIRECTORY)
    
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

if __name__ == "__main__":
    main()