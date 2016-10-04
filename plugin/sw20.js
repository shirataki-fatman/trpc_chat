var commandTable = [
    [/k(\d+)(@(\d+))?\s?/, RollPowerDice]
];

for (var i in commandTable) {
    var matchData = prop.Message.match(commandTable[i][0])
    if (matchData) {
        commandTable[i][1](matchData);
    }
}

function RollPowerDice(matchData) {
    var powerTable = {
         0: [-1,  0,  0,  0,  1,  2,  2,  3,  3,  4,  4],
        10: [-1,  1,  1,  2,  3,  3,  4,  5,  5,  6,  7],
        20: [-1,  1,  2,  3,  4,  5,  6,  7,  8,  9, 10],
        30: [-1,  2,  4,  4,  6,  7,  8,  9, 10, 10, 10],
        40: [-1,  4,  5,  6,  7,  9, 10, 11, 11, 12, 13],
        50: [-1,  4,  6,  8, 10, 10, 12, 12, 13, 15, 15],
    };

    var power = matchData[1];
    var critical = matchData[3] || 10;
    if (!powerTable[power] || critical < 2) {
        return;
    }
    var total = 0;
    var numberStrings = [];
    var totals = [];
    var damages = [];
    var totalDamage = 0;
    do {
        var diceResult = prop.RollDice(2, 1, 6);
        total = parseInt(diceResult[0]) + parseInt(diceResult[1]);
        numberStrings.push(diceResult.join(","));
        totals.push(total);
        if (powerTable[power][total - 1] > 0) {
            damages.push(powerTable[power][total - 1]);
            totalDamage += powerTable[power][total - 1];
        }
    } while(total >= critical && powerTable[power][total - 1]);

    var resultString = "2D:[" + numberStrings.join(" ") + "]=" + totals.join(",") + " -> " + damages.join(",");
    if (numberStrings.length > 1) {
        resultString += " -> " + (numberStrings.length - 1) + "回転";
    }
    resultString += " -> " + totalDamage;
    prop.Result = resultString;
}