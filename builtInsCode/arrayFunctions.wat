(module
    (memory (export "memory") 1)
    (type $0 (func (param i32) (result i32)))
    (func $len (type $0) (param $arrayPointer i32) (result i32)
        (i32.load (local.get $arrayPointer))
    )
)
    
