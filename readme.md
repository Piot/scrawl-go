### scrawl-go

Interface Description Language.

##### Usage

```go
definition, definitionErr := scrawl.ParseString("type Wheel\n  angle int\n")
```


##### Interface file
Each indentation step must be defined with exactly two space characters. The basic types is up to your implementation to define.

###### Example
```
type Transform
  position Position
  rotation Rotation
    
component Wheel
  angle int32
    
component Body
  transform Transform
    
entity Cycle
  body Body
  front Wheel
  back Wheel   
```
