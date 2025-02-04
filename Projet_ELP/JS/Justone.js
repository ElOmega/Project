function collectClues(players){
    const clues=[];
    for (const player of players){
        while (true){
            const clue=prompt(player+", qu'elle est ton mot clé");
            console.log(player + "a choisi le mot :"+ clue);
            if (clues.includes(clue)){
                console.log("Le clue est déjà pris");
            }
            else{
                clues.push(clue);
                break;
            }
        }
    }
    return clues;
}

function startGame(){
    const words=["Arbre","Pomme","Foot","Chien"];
    const words_to_guess=words[Math.floor(Math.random()*words.length)];
    console.log("Le mot a faire deviner est " + words_to_guess);
    return words_to_guess
}

function display(clues,words_to_guess){
    console.log("ton mot à définir est composer de" + words_to_guess.length + "lettres");
    console.log("Voici les mots clé:");
    for (const mot of clues){
        console.log(mot);
    }
}

function checkGuess(guess_word,word_to_guess){

}

function playGame(){
    const players =[];
    while (true){
        const nom = prompt("Comment s'appelle le nom du joueur");
        players.push(nom);
        const verif = prompt ("c'est tout les joueurs ?");
        if (verif == "oui"){
            break;
        }
    }
    const word_to_guess = startGame();
    const clues = collectClues(players);
    display(clues,word_to_guess);
    console.log("Tu as 5 essaies pour le reussir");
    const number_of_try=5
    for(let i=0 ; i<number_of_try ; i++){
        const guess_word=prompt("Qu'elle mot tu pense que c'est ?");
        if (guess_word===word_to_guess){
            console.log("Félécitation tu as le bon mot")
            return
        } else{
            console.log("C'est le mauvais mot tu as plus que " + (number_of_try-i))
        }
    }
    console.log("Tu as fait" + number_of_try + "error. Voici le mot qu'il fallait deviner" + word_to_guess + "!")
}

playGame()