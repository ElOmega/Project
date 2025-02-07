const readline = require("readline");

const rl = readline.createInterface({
  input: process.stdin,
  output: process.stdout,
});

function askQuestion(query) {
  return new Promise((resolve) => rl.question(query, resolve));
}

async function collectClues(players) {
  const clues = [];
  for (const player of players) {
    while (true) {
      const clue = await askQuestion(player + ", quel est ton mot clé ? ");
      console.log({player} + " a choisi le mot : " + clue);
      if (clues.includes(clue)) {
        console.log("Le mot clé est déjà pris !");
      } else {
        clues.push(clue);
        break;
      }
    }
  }
  return clues;
}

function startGame() {
  const words = ["Arbre", "Pomme", "Foot", "Chien"];
  const word_guess = words[Math.floor(Math.random() * words.length)];
  console.log("Le mot à faire deviner est : " + word_guess);
  return word_guess;
}

function display(clues, word_guess) {
  console.log("Le mot à définir est composé de " + word_guess.length + " lettres.");
  console.log("Voici les mots clés : ");
  for (const clue of clues) {
    console.log(clue);
  }
}

async function playGame() {
  const players = [];
  while (true) {
    const name = await askQuestion("Comment s'appelle le joueur ? ");
    players.push(name);
    const fin = await askQuestion("C'est tout pour les joueurs ? (oui/non) ");
    if (fin.toLowerCase() === "oui") {
      break;
    }
  }

  const word_guess = startGame();
  const clues = await collectClues(players);
  console.log('\x1Bc');
  display(clues, word_guess);

  console.log("Tu as 5 essais pour trouver le mot.");
  const numberOfTries = 5;
  for (let i = 0; i < numberOfTries; i++) {
    const word_guessed = await askQuestion("Quel mot penses-tu que c'est ? ");
    if (word_guessed.toLowerCase() === word_guess.toLowerCase()) {
      console.log('\x1BC')
      console.log("Félicitations, tu as trouvé le bon mot !");
      rl.close();
      return;
    } else {
      console.log(
        "Mauvaise réponse. Il te reste " + (numberOfTries - (i + 1)) + " essais."
      );
    }
  }
  console.log('\x1BC')
  console.log(
    "Tu as épuisé tes" + numberOfTries + " essais. Le mot à deviner était : ${wordToGuess}."
  );
  rl.close();
}

playGame();
