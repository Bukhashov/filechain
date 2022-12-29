import os
import logging
import face_recognition
import grpc
import face_pb2, face_pb2_grpc
from concurrent import futures

### tmp file ға сақтауды жасауды істеу керек
def get_pathfile(filename, extension):
    return f'{filename}{extension}'

class Face(face_pb2_grpc.FaceServicer):
    def Find(self, request_iterator, context):
        data = bytearray()
        filepath = "" 
       
        for request in request_iterator:
            if request.metadata.filename and request.metadata.extension:
                filepath = get_pathfile(request.metadata.filename, request.metadata.extension)
                continue
            data.extend(request.image)
        with open(filepath, 'wb') as f:
            f.write(data)
            f.close()
        
        face_img = face_recognition.load_image_file(filepath)
        faceArr = face_recognition.face_locations(face_img)
        
        # os.remove(filepath)
        print(faceArr)
        return face_pb2.FindRespons(total=len(faceArr))

    def Comparison(self, request_iterator, context):
        data_original = bytearray()
        data_forCheck = bytearray()
        filepath_original = ""
        filepath_forCheck = ""
        res = False

        for request in request_iterator:
            if request.originalMetadata.filename and request.originalMetadata.extension:
                filepath_original = get_pathfile(request.originalMetadata.filename, request.originalMetadata.extension)
                continue
            if request.forCheck.filename and request.forCheck.extension:
                filepath_forCheck = get_pathfile(request.forCheck.filename, request.forCheck.extension)
                continue
            data_original.extend(request.originalImage)
            data_forCheck.extend(request.forCheckImage)
        # print(data_original)
        # print(data_forCheck)
        with open("./org/"+filepath_original, 'wb') as f:
            f.write(data_original)
            f.close()
        with open("./tmp/"+filepath_forCheck, 'wb') as f:
            f.write(data_forCheck)
            f.close()

        img_original = face_recognition.load_image_file("./org/"+filepath_original)
        img_forCheck = face_recognition.load_image_file("./tmp/"+filepath_forCheck)

        encode_original = face_recognition.face_encodings(img_original)[0]
        encode_forCheck = face_recognition.face_encodings(img_forCheck)[0]

        result = face_recognition.compare_faces([encode_original], encode_forCheck)

        if result[0] == True:
            res = True
        elif result[0] == False:
            res = False
        else:
            print(result[0])
        
        os.remove("./org/"+filepath_original)
        os.remove("./tmp/"+filepath_forCheck)
        
        return face_pb2.ComparisonRespons(coincidences=res)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    face_pb2_grpc.add_FaceServicer_to_server(Face(), server)
    server.add_insecure_port('[::]:5050')
    server.start()
    server.wait_for_termination()

if __name__ == "__main__":
    logging.basicConfig()
    serve()
