module Main exposing (main)

import Browser
import Html exposing (Html, div, input, button, text)
import Html.Attributes exposing (placeholder, value)
import Html.Events exposing (onClick, onInput)
import Svg exposing (Svg, svg, line)
import Svg.Attributes exposing (x1, x2, y1, y2, stroke, strokeWidth, width, height)
import Basics exposing (degrees, cos, sin)
import Types exposing (Instruction(..), Program)


-- MODEL

type alias Position =
    { x : Float, y : Float }

type alias Model =
    { userInput : String  -- The raw text input
    , program : Program   -- Parsed program from the user input
    , position : Position
    , angle : Float
    }

init : Model
init =
    { userInput = ""
    , program = []
    , position = { x = 250, y = 250 }
    , angle = 0
    }


-- DRAWING FUNCTIONS

calculateNewPosition : Float -> Float -> Float -> Float -> (Float, Float)
calculateNewPosition x y angle distance =
    let
        radianAngle = degrees angle
        newX = x + distance * cos radianAngle
        newY = y + distance * sin radianAngle
    in
    (newX, newY)

createSvgLine : Float -> Float -> Float -> Float -> Svg msg
createSvgLine xStart yStart xEnd yEnd =
    line
        [ x1 (String.fromFloat xStart)
        , y1 (String.fromFloat yStart)
        , x2 (String.fromFloat xEnd)
        , y2 (String.fromFloat yEnd)
        , stroke "black"
        , strokeWidth "3"
        ]
        []

drawProgram : Program -> Position -> Float -> List (Svg msg) -> List (Svg msg)
drawProgram instructions position angle svgElements =
    case instructions of
        [] ->
            svgElements

        instruction :: rest ->
            case instruction of
                Forward distance ->
                    let
                        (newX, newY) = calculateNewPosition position.x position.y angle distance
                        newLine = createSvgLine position.x position.y newX newY
                        newPosition = { x = newX, y = newY }
                    in
                    drawProgram rest newPosition angle (newLine :: svgElements)

                Left turnAngle ->
                    drawProgram rest position (angle - turnAngle) svgElements

                Right turnAngle ->
                    drawProgram rest position (angle + turnAngle) svgElements

                Repeat n repeatedInstructions ->
                    let
                        repeatedProgram = List.concat (List.repeat n repeatedInstructions)
                    in
                    drawProgram (repeatedProgram ++ rest) position angle svgElements


-- VIEW FUNCTION

view : Model -> Html Msg
view model =
    div []
        [ input
            [ placeholder "Enter commands (e.g., F100 R90 F100)"
            , value model.userInput
            , onInput UserInputChanged
            ]
            []
        , button [ onClick ParseCommands ] [ text "Draw" ]
        , svg [ width "500", height "500" ] (List.reverse (drawProgram model.program model.position model.angle []))
        ]


-- USER INPUT & PARSING

type Msg
    = UserInputChanged String
    | ParseCommands

update : Msg -> Model -> Model
update msg model =
    case msg of
        UserInputChanged newInput ->
            { model | userInput = newInput }

        ParseCommands ->
            { model | program = parseCommands model.userInput }


parseCommands : String -> Program
parseCommands input =
    let
        tokens = String.words input
    in
    parseTokens tokens []


parseTokens : List String -> Program -> Program
parseTokens tokens acc =
    case tokens of
        [] ->
            List.reverse acc

        "F" :: n :: rest ->
            case String.toInt n of
                Just distance -> parseTokens rest (Forward distance :: acc)
                Nothing -> parseTokens rest acc

        "L" :: n :: rest ->
            case String.toInt n of
                Just angle -> parseTokens rest (Left angle :: acc)
                Nothing -> parseTokens rest acc

        "R" :: n :: rest ->
            case String.toInt n of
                Just angle -> parseTokens rest (Right angle :: acc)
                Nothing -> parseTokens rest acc

        "REP" :: n :: rest ->
            case String.toInt n of
                Just times ->
                    let
                        (repeatedInstructions, remaining) = extractRepeatBlock rest
                    in
                    parseTokens remaining (Repeat times repeatedInstructions :: acc)
                Nothing ->
                    parseTokens rest acc

        _ :: rest ->
            parseTokens rest acc


extractRepeatBlock : List String -> (Program, List String)
extractRepeatBlock tokens =
    let
        extract innerTokens acc =
            case innerTokens of
                "END" :: remaining ->
                    (List.reverse acc, remaining)

                [] ->
                    (List.reverse acc, [])

                token :: rest ->
                    let
                        parsed = parseTokens [token] []
                    in
                    extract rest (acc ++ parsed)
    in
    extract tokens []


-- MAIN

main =
    Browser.element { init = \_ -> init, update = update, view = view, subscriptions = \_ -> Sub.none }
