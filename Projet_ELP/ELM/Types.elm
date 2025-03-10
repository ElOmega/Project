module Types exposing (Instruction(..), Program)

type Instruction
    = Forward Int
    | Left Int
    | Right Int
    | Repeat Int (List Instruction)

type alias Program =
    List Instruction
