(module
  (memory (export "memory") 1)
  (type $0 (func (param i32) (result i32)))
  (type $1 (func (param i32)))
  (type $2 (func (param i32) (param i32) (result i32)))
  (type $3 (func (param i32) (param i32) (param i32) (result i32)))
 
  ;; Gets num bytes and returns pointer to chunk with minimum num bytes lenght
  (func $allocate (type $0) (param $numBytes i32) (result i32)
    (local $nextChunkPos i32)
    (local $nextChunkLen i32)
    (local $orginalChunkLen i32)
    (local $memoryPointer i32)
    (local.set $memoryPointer (i32.const 0))

    (block $0
      (loop $1
        ;; TODO: cheack if memory pointer is larger than curent meomory size

        (i32.eq (i32.load (local.get $memoryPointer)) (i32.const 0))
        (if  ;; If the current chunk is unused
          (then
            (i32.eq (i32.load (i32.add (local.get $memoryPointer) (i32.const 4))) (i32.const 0))
            (if  ;; If the position that gives the lenght of the chunk of meomory is zero the memory is unused
              (then 
                (br $0)
              ) 
            )

            (i32.le_u (local.get $numBytes) (i32.load (i32.add (local.get $memoryPointer) (i32.const 4)))) 
            (if ;; If number of bytes to be allocated is grater than the lenght of the chunk
              (then 
                (br $0)
              ) 
            )
          )
        ) 
        
        (local.set $memoryPointer 
          (i32.add 
            (i32.add
              (local.get $memoryPointer) 
              (i32.load (i32.add (local.get $memoryPointer) (i32.const 4))))
            (i32.const 8) ;; the 4 bytes that store if the chunk is used + the 4 bytes that store the chunk lenght  
          ) 
        ) ;; memoryPointer += chunkLen + chunkLenByte + usedByte

        (br $1) ;; Loop
      )
    )

    (i32.store (local.get $memoryPointer) (i32.const 1)) ;; Set chunk to used
    (local.set $orginalChunkLen (i32.load (i32.add (local.get $memoryPointer) (i32.const 4)))) ;; Get the size of the chunk

    (if (i32.eq (local.get $orginalChunkLen) (i32.const 0)) ;; If the size of the chunk is zero there are no chunks later in memory so the allocater can set the chunk lenght to numBytes without spliting the chunk as long as the memoryPointer + numBytes + 2 is not longer then the size of memory
      (then 
        ;; TODO: cheack if the allocater needs to increasse the memory size

        (i32.store (i32.add (local.get $memoryPointer) (i32.const 4)) (local.get $numBytes))

        (i32.add (local.get $memoryPointer ) (i32.const 8))
        return 
      ) 
    )

    (if (i32.gt_u (i32.add (local.get $numBytes) (i32.const 24)) (i32.load (i32.add (local.get $memoryPointer) (i32.const 4))) ) ;; if numBytes + 24 >  chunkLen (No point in spliting the chunk if the split chunk is so small)
      (then 
        (i32.add (local.get $memoryPointer ) (i32.const 8))
        return 
      )
    )

    ;;Split the chunks 
    (i32.store (i32.add (local.get $memoryPointer) (i32.const 4)) (local.get $numBytes)) ;; Set the first chunk len to numBytes

    (local.set $nextChunkPos (i32.add (i32.add (local.get $numBytes) (local.get $memoryPointer)) (i32.const 8)))
    (i32.store (local.get $nextChunkPos) (i32.const 0)) ;; Set next chunk to unused
    (i32.store (i32.add (local.get $nextChunkPos) (i32.const 4)) (i32.sub (i32.sub (local.get $orginalChunkLen) (local.get $numBytes)) (i32.const 8))) ;; nextChunkLen = orginalChunkLen - numBytes - 8

    (i32.add (local.get $memoryPointer ) (i32.const 8))
    return 
  )

  (func $deAllocate (type $1) (param $chunkPointer i32)
    (i32.store (i32.sub (local.get $chunkPointer) (i32.const 8)) (i32.const 0)) ;; Set chunk to unused
    ;; TODO: defragmantion
  )

  (func $array (type $2) (param $numElements i32) (param $elementSize i32) (result i32)
    (local $arrayPointer i32)
    (local $i i32)
    (local $endLoop i32)

    (i32.add (i32.mul (local.get $numElements) (local.get $elementSize)) (i32.const 4)) ;; adding four to store lenght
    (call $allocate)
    (local.set $arrayPointer)

    (i32.store (local.get $arrayPointer) (local.get $numElements))
    (local.set $i (i32.add (local.get $arrayPointer) (i32.const 4)))
    (local.set $endLoop (i32.add (i32.mul (local.get $elementSize) (local.get $numElements)) (i32.add (local.get $arrayPointer) (i32.const 4))))

    ;; set all bytes to 0
    (block $0
      (loop $1
        (i32.eq (local.get $i) (local.get $endLoop))
        (if 
          (then
            (br $0)
          )
        )

        (i32.store8 (local.get $i) (i32.const 0))
        (local.set $i (i32.add (local.get $i) (i32.const 1)))
        
        (br $1)
      )
    )

    (local.get $arrayPointer)
    return
  )

  (func $take (type $3) (param $numToTake i32) (param $arrayPointer i32) (param $elementSize i32) (result i32)
    (local $newArray i32)
    (local $newArrayFirstPos i32)
    (local $arrayPointerFirstPos i32)
    (local $i i32)
    (local $numBytesToTake i32)

    (if (i32.gt_u (local.get $numToTake) (i32.load (local.get $arrayPointer))) (then (unreachable)))
    
    (local.set $newArray (call $array (local.get $numToTake) (local.get $elementSize)))
    (local.set $newArrayFirstPos (i32.add (local.get $newArray) (i32.const 4)))
    (local.set $arrayPointerFirstPos (i32.add (local.get $arrayPointer) (i32.const 4)))
    (local.set $numBytesToTake (i32.mul (local.get $numToTake) (local.get $elementSize)))

    (loop $copy
      (i32.store8
        (i32.add (local.get $newArrayFirstPos) (local.get $i))
        (i32.load8_u (i32.add (local.get $arrayPointerFirstPos) (local.get $i)))
      )
      
      (local.set $i (i32.add (local.get $i) (i32.const 1)))
      (br_if $copy (i32.lt_u (local.get $i) (local.get $numBytesToTake)))
    )

    (local.get $newArray)
  )

  (func $tail (type $2) (param $arrayPointer i32) (param $elementSize i32) (result i32)
    (local $oldArraySize i32)
    (local $newArraySize i32)
    (local $newArray i32)
    (local $newArrayFirstPos i32)
    (local $oldArrayFirstPos i32)
    (local $i i32)
    (local $numBytesToTake i32)

    (local.set $oldArraySize (i32.load (local.get $arrayPointer)))
    (local.set $newArraySize (i32.sub (local.get $oldArraySize) (i32.const 1)))

    (if (i32.le_u (local.get $oldArraySize) (i32.const 1)) 
      (then (return (call $array (i32.const 0) (local.get $elementSize))))
    )
    
    (local.set $newArray (call $array (local.get $newArraySize) (local.get $elementSize)))
    (local.set $newArrayFirstPos (i32.add (local.get $newArray) (i32.const 4))) (; Adding four to skip length ;)
    (local.set $oldArrayFirstPos (i32.add (i32.add (local.get $arrayPointer) (i32.const 4)) (local.get $elementSize))) (; Old array first pos storing pointer to the second element in the array ;)
    (local.set $numBytesToTake (i32.mul (local.get $oldArraySize) (local.get $elementSize))) 

    (loop $copy
      (i32.store8
        (i32.add (local.get $newArrayFirstPos) (local.get $i))
        (i32.load8_u (i32.add (local.get $oldArrayFirstPos) (local.get $i)))
      )
      
      (local.set $i (i32.add (local.get $i) (i32.const 1)))
      (br_if $copy (i32.lt_u (local.get $i) (local.get $numBytesToTake))) (; Loop if i is less then numBytes ;)
    )

    (local.get $newArray)
  )
)

(;
The first 4 bytes is used to store if the chunk is used

[chunk is used i.32, chunk lenght i.32, elem1, elem2, elem3, ... chunk is used, chunk lenght, elem1, elem2, elem3]

[0 10 0 0 0 0 0 0 0 0 0 0 1 2 0 0]

orginalChunkLen - numBytes - 8

;)