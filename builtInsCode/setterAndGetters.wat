(module
    (memory (export "memory") 1)
    (type $0 (func (param i32) (param i32) (result i32)))
    (type $1 (func (param i32) (param i32) (result f32)))
    (type $2 (func (param i32) (param i32) (param i32) (result i32)))
    (type $3 (func (param i32) (param i32) (param f32) (result i32)))

    (func $i32get (type $0) (param $arrayPointer i32) (param $index i32) (result i32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (i32.load (i32.add (local.get $arrayPointer) (i32.mul (local.get $index) (i32.const 4))))
    )

    (func $i8get (type $0) (param $arrayPointer i32) (param $index i32) (result i32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (i32.load8_u (i32.add (local.get $arrayPointer) (local.get $index)))
    )

    (func $f32get (type $1) (param $arrayPointer i32) (param $index i32) (result f32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (f32.load (i32.add (local.get $arrayPointer) (i32.mul (local.get $index) (i32.const 4))))
    )

    (func $i32set (type $2) (param $arrayPointer i32) (param $index i32) (param $elementValue i32) (result i32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (i32.store (i32.add (local.get $arrayPointer) (i32.mul (local.get $index) (i32.const 4))) (local.get $elementValue))
        (local.get $arrayPointer)
    )

    (func $i8set (type $2) (param $arrayPointer i32) (param $index i32) (param $elementValue i32) (result i32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (i32.store8 (i32.add (local.get $arrayPointer) (local.get $index)) (local.get $elementValue))
        (local.get $arrayPointer)
    )

    (func $f32set (type $3) (param $arrayPointer i32) (param $index i32) (param $elementValue f32) (result i32)
        (if (i32.ge_u (local.get $index) (i32.load (local.get $arrayPointer))) (then (unreachable)))

        (local.set $arrayPointer (i32.add (local.get $arrayPointer) (i32.const 4)))
        (f32.store (i32.add (local.get $arrayPointer) (i32.mul (local.get $index) (i32.const 4))) (local.get $elementValue))
        (local.get $arrayPointer)
    ) 
)