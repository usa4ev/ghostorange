## Ghost Orange

This is basicly a storage that can store various types of data like **`text data, binary data, credentials and bank card info`**. It provides simple tui that access server via http.

![general scheme](./assets/data_flow_scheme.svg)

### Features:
Service implements server-side encryption for credentials datatype. 

Access to bank cards data requres authorization via CVV-code input. The code is not stored openly.

To run server with default config you can use MakeFile.
On linux run:
```
make run-srv-linux
``` 

On windows run:
```
make run-srv-windows
``` 

On mac run:
```
make run-srv-darwin
``` 

Same goes for client.

On linux run:
```
make run-tui-linux
``` 

On windows run:
```
make run-tui-windows
``` 

On mac run:
```
make run-tui-darwin
``` 

You can also run tests with:
```
make test
``` 