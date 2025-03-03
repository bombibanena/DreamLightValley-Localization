import grpc
from concurrent import futures
from proto import translation_pb2
from proto import translation_pb2_grpc
from translator import translate

class TranslationService(translation_pb2_grpc.TranslationServiceServicer):
    def Translate(self, request, context):
        text = request.text
        
        data = translate(text)
        
        return translation_pb2.TranslationResponse(data=data)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    translation_pb2_grpc.add_TranslationServiceServicer_to_server(TranslationService(), server)
    server.add_insecure_port('[::]:50051')
    print("Server started on port 50051...")
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    serve()
