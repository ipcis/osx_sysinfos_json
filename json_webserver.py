import json
from http.server import BaseHTTPRequestHandler, HTTPServer

#$ curl -X POST -H "Content-Type: application/json" -d '{"name": "Alice", "age": 30}' http://localhost:8080

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
    def do_POST(self):
        content_length = int(self.headers['Content-Length'])
        post_data = self.rfile.read(content_length)
        json_data = json.loads(post_data.decode('utf-8'))
        print(json_data)
        self.send_response(200)
        self.end_headers()

if __name__ == '__main__':
    web_server = HTTPServer(('localhost', 8080), SimpleHTTPRequestHandler)
    web_server.serve_forever()
