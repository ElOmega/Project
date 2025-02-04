module Main exposing (main)

import Browser
import Html exposing (Html, div, text, input, br, a, button)
import Html.Attributes exposing (style, placeholder, value, href)
import Html.Events exposing (onInput, onClick)


-- MODEL

type alias Model =
    { userInput : String
    , errorMessage : String
    }

init : Model
init =
    { userInput = ""
    , errorMessage = ""
    }


-- UPDATE

type Msg
    = UserInputChanged String
    | DessinerClicked

update : Msg -> Model -> Model
update msg model =
    case msg of
        -- Mise à jour de la saisie utilisateur sans analyser à chaque frappe
        UserInputChanged newInput ->
            { model | userInput = newInput }

        -- Lors du clic sur le bouton « Dessiner », on analyse l'entrée
        DessinerClicked ->
            if isValidInput model.userInput then
                { model | errorMessage = "" }
            else
                { model | errorMessage = "Mauvaise syntaxe, besoin d'aide?" }


-- VALIDATION DE L'ENTRÉE

isValidInput : String -> Bool
isValidInput input = 
    let
        pattern =
            "^\\[.*(?:Forward|Left|Right|Repeat)\\s+\\d+(?:\\s*,\\s*(?:Forward|Left|Right|Repeat)\\s+\\d+)*.*\\]$"

        regex =
            Regex.fromString pattern |> Maybe.withDefault Regex.never
    in
    Regex.contains regex input


-- VIEW

view : Model -> Html Msg
view model =
    div
        [ style "display" "flex"
        , style "flex-direction" "column"
        , style "justify-content" "center"
        , style "align-items" "center"
        , style "height" "100vh"
        , style "font-family" "Calibri, sans-serif"
        ]
        [ div [ style "font-size" "48px" ] [ text "LHOMME, BORDES, BLANZAT: Projet ELM" ]
        , div [ style "font-size" "25px" ] [ text "Entrez les commandes ci-dessous:" ]
        , br [] []
        , input 
            [ placeholder "Exemple: [Forward 100, Repeat 180[Left 100]]"
            , value model.userInput
            , onInput UserInputChanged
            , style "width" "600px"
            ]
            []
        , button 
            [ onClick DessinerClicked
            , style "margin-top" "10px"
            , style "font-size" "20px"
            ]
            [ text "Dessiner" ]
        , div [ style "margin-top" "10px", style "font-size" "20px" ] [ text model.userInput ]
        , div [ style "margin-top" "10px", style "font-size" "20px", style "color" "red" ]
            [ a 
                [ href "https://perso.liris.cnrs.fr/tristan.roussillon/TcTurtle/project.md.html"
                , style "color" "blue"
                , style "text-decoration" "underline"
                ]
                [ text model.errorMessage ]
            ]
        , div 
            [ style "width" "400px"
            , style "height" "100px"
            , style "min-height" "400px"
            , style "border" "1px solid black"
            , style "margin-top" "20px"
            , style "display" "block"
            , style "box-sizing" "border-box"
            ]
            [ text " " ]
        ]


-- MAIN

main : Program () Model Msg
main =
    Browser.sandbox { init = init, update = update, view = view }