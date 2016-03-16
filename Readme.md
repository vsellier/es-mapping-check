# ES index checker

Found index and document types with field with dot in proprty names. 

## Usage

```
go run index-checker.go
```

Use a local elasticsearch instance on default port by default 

Example of output :
```
Loading index list....
378/426 : index-name property system.property.name not compatible


Indices to fix : index-name
```