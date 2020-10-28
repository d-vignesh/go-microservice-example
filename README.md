This is simple application that demonstrates building microservices with golang. The application shows the various coffee products available from a vender and consists of three services namely product-frontend, product-api, currency. 

The product-frontend service is build with react js. The service makes the REST calls to the product-api service to fetch details of the list of coffee products.

The product-api is a service that process the request from the frontend service and responds back with the required data in REST format. The product-api gets the list of available coffee products from the data storage(which is this case is a simple go slice). The product-api also make a call to the currency service to get the exchange rate if the currency requested from the frontend is different from default. The communication between the product-api and currency service is carried out via gRPC. 

The currency service is responsible to making the query for the updated exchange rates for different currency and serve it to the product-api when requested. 

To run the application, you need to install go and npm intalled and clone the repo and perform the below steps,
1. Run the currency service using below steps:
    cd currency/
    go run main.go (starts the servemux at localhost:9090)
2. Run the product-api service using below steps:
    cd product-api/
    go run main.go (starts the servemux at localhost:9092)
3. Run the product-frontend service using below steps:
    cd product-frontend/
    npm start (startes the react js app at localhost:3000)

This application was build with reference to an excellent youtube vedio series from Nic Jackon which gives a clean explanation of these concepts: https://www.youtube.com/playlist?list=PLmD8u-IFdreyh6EUfevBcbiuCKzFk0EW_