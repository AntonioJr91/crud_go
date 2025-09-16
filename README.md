#Technologies  
Backend: Go  
Frontend: JavaScript / HTML  
Structure: directories backend/ and frontend  

#Features  
Create, read, update, and delete items (CRUD)  
REST API on the backend  
Simple frontend interface to consume the API  

#How to run  
Clone the repository  
git clone https://github.com/AntonioJr91/crud_go.git  

cd crud_go  

#Start the backend  
cd backend  
go run .  

#Unit tests  
cd backend  
go test -v ./...  

#Docker  
frontend local port 3000  
backend local port 8080  
mariadb local port 3306  

start the containers  
docker compose up --build
