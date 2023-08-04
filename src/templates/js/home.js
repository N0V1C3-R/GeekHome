function getCookie(name) {
    let cookieArr = document.cookie.split(";");

    for (let i = 0; i < cookieArr.length; i++) {
        let cookiePair = cookieArr[i].trim().split("=");
        let cookieName = cookiePair[0];
        let cookieValue = cookiePair[1];

        if (cookieName === name) {
            return decodeURIComponent(cookieValue);
        }
    }
    return null;
}

function includesIgnoreCase(stringList, string) {
    const stringListUpper = stringList.map(str => str.toUpperCase());
    return stringListUpper.includes(string.toUpperCase())
}

async function postRequestServer(url, body) {
    return await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: body
    })
}

const commandInput = document.querySelector('.command__input');

commandInput.addEventListener("keydown", (event) => {
    if (event.key === "Enter") {
        commandInput.setAttribute("disabled", "disabled");
        const _ = commandListener(event);
        event.preventDefault();
    }
});

async function addNewTerminal(username, ip) {
    const terminal = document.querySelector(".terminal");
    terminal.classList.add("terminal");
    const newCommand = document.createElement("div");
    newCommand.classList.add("command");
    const newUserInfo = document.createElement("span");
    newUserInfo.classList.add("user__info");
    newUserInfo.textContent = username + "@" + ip + "\u00A0/\u00A0$\u00A0";
    const newCommandInput = document.createElement("input");
    newCommandInput.classList.add("command__input");
    newCommandInput.type = "text";
    newCommandInput.autofocus = true;
    newCommandInput.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            newCommandInput.setAttribute("disabled", "disabled");
            commandListener(event);
            event.preventDefault();
        }
    });
    newCommand.appendChild(newUserInfo);
    newCommand.appendChild(newCommandInput);
    terminal.appendChild(newCommand);
    newCommandInput.focus();
}

class Command {
    constructor(commandName) {
        this.commandName = commandName;
        this.aliasList = [];
        this.params = {};
        this.resp = undefined;
    };

    alias(aliasName) {
        this.aliasList.add(aliasName);
    };

    async execute(inputText) {
        var body = JSON.stringify({"stdin": inputText})
        // this.resp = "ERROR: [command] command that cannot be parsed".replace("[command]", "[" + this.commandName + "]");
        var responseObj = await this.postRequestServer("/command", body)
        if (responseObj.redirected) {
            console.log(responseObj.url)
            window.location.href = responseObj.url;
        }
        var response = await responseObj.json()
        this.resp = response["response"]
        this.printResponse()
    }

    async postRequestServer(url, body) {
        return await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: body
        })
    };

    printResponse() {
        if (this.resp !== "undefined") {
            const currentTerminalElement = document.querySelector(".terminal")
            const result = document.createElement("div");
            result.classList.add("result");
            result.innerHTML = this.resp;
            currentTerminalElement.appendChild(result);
        }
    }
}

class ClearCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["clear", "cls"];
    }

    async execute(inputList) {
        if (inputList.length !== 0) {
            this.resp = "ERROR: [command] command that cannot be parsed".replace("[command]", "[" + this.commandName + " " + inputList.join(" ") + "]");
            this.printResponse();
        } else {
            const fatherDiv = document.querySelector(".terminal");
            fatherDiv.innerHTML = "";
        }
    }
}

class OpenCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["open"];
    }

    async execute(inputList) {
        if (inputList.length === 0) {
            this.resp = "ERROR: [command] command that cannot be parsed".replace("[command]", "[" + this.commandName + "]");
        } else {
            const protocolHeaders = ["https://", "http://", "file:///", "ftp://"]
            switch (inputList.length) {
                case 1:
                    if (protocolHeaders.some(protocolHeader => inputList[0].toLowerCase().startsWith(protocolHeader))) {
                        window.open(inputList[0], "_blank");
                    } else {
                        this.resp = "ERROR: Temporarily unsupported protocols, or protocols missing. HTTP, HTTPS, FILE, FTP are supported protocols for this application."
                    }
                    break;
                case 2:
                    if (protocolHeaders.some(protocolHeader => inputList[1].toLowerCase().startsWith(protocolHeader))) {
                        if (inputList[0].toUpperCase() === "-R") {
                            window.location.href = inputList[0];
                        } else if (inputList[0].toUpperCase() === "-N") {
                            window.open(inputList[0], "_blank");
                        } else {
                            this.resp = "Error: Illegal mode\n"
                                + "open -N :(default) Open new window\n"
                                + "\t-R :Current window open link\n"
                        }
                    } else {
                        this.resp = "ERROR: Temporarily unsupported protocols, or protocols missing. HTTP, HTTPS, FILE, FTP are supported protocols for this application."
                    }
                    break;
                default:
                    this.resp = "ERROR: Commands that cannot be parsed."
                    break;
            }
        }
        if (typeof this.resp !== "undefined") {
            this.printResponse();
        }
    }
}

class LoginCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["login"]
    }

    async execute(inputList) {
        if (inputList.length !== 0) {
            this.resp = "ERROR: [command] command that cannot be parsed".replace("[command]", "[" + this.commandName + " " + inputList.join(" ") + "]");
        } else {
            window.location.href = "/login";
            return
        }
        if (typeof this.resp !== "undefined") {
            this.printResponse();
        }
    }
}

class ExitCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["logout", "exit"];
    }

    async execute(inputList) {
        const response = await postRequestServer("/logout", null);
        const data = await response.json();
        this.resp = data["response"];
        if (this.commandName.toLowerCase() === "logout" && this.resp === true) {
            this.resp = "Logout Successful!"
            this.printResponse();
        } else if (this.commandName.toLowerCase() === "logout" && this.resp === false) {
            this.resp = "ERROR: No login info."
            this.printResponse();
        } else if (this.commandName.toLowerCase() === "exit"  && this.resp === true) {
            this.resp = "Logout Successful!"
            this.printResponse();
        } else if (this.commandName.toLowerCase() === "exit"  && this.resp === false) {
            window.open('', '_self').close();
        }

    }
}

class RegisterCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["register", "signup"]
    }

    async execute(inputList) {
        if (inputList.length !== 0) {
            this.resp = "ERROR: [command] command that cannot be parsed".replace("[command]", "[" + this.commandName + " " + inputList.join(" ") + "]");
        } else {
            window.location.href = "/register";
            return
        }
        if (typeof this.resp !== "undefined") {
            this.printResponse();
        }
    }
}

class GoogleCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["google", "go"]
    }

    async execute(inputList) {
        if (inputList.length === 0) {
            window.open("https://www.google.com")
        } else {
            inputList = inputList.map(str => str.replace("%", "%25"))
            inputList = inputList.map(str => str.replace("+", "%2B"))
            const search = inputList.join("%20")
            window.open( "https://www.google.com/search?q=" + search, "_blank");
        }
    }
}

class BingCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["bing"]
    }

    async execute(inputList) {
        if (inputList.length === 0) {
            window.open("https://www.bing.com")
        } else {
            inputList = inputList.map(str => str.replace("%", "%25"))
            inputList = inputList.map(str => str.replace("+", "%2B"))
            const search = inputList.join("%20")
            window.open( "https://www.bing.com/search?q=" + search, "_blank");
        }
    }
}

class GithubCommand extends Command {
    constructor(commandName) {
        super(commandName);
        this.aliasList = ["github"]
    }

    async execute(inputList) {
        if (inputList.length === 0) {
            window.open("https://github.com/")
        } else {
            inputList = inputList.map(str => str.replace("%", "%25"))
            inputList = inputList.map(str => str.replace("+", "%2B"))
            const search = inputList.join("%20")
            window.open( "https://github.com/search?q=" + search, "_blank");
        }
    }
}

async function commandListener(event) {
    const inputText = event.target.value.trim();
    const textList = inputText.split(' ')
    const command = textList[0];

    const commandClasses = [
        ClearCommand, OpenCommand, LoginCommand, ExitCommand, RegisterCommand, GoogleCommand, BingCommand, GithubCommand
    ];
    let commandClass;
    for (const cls of commandClasses) {
        const commandInstance = new cls(command);
        if (includesIgnoreCase(commandInstance.aliasList, command)) {
            commandClass = commandInstance;
        }
    }
    if (typeof commandClass === "undefined") {
        commandClass = new Command(command);
        await commandClass.execute(inputText)
    } else {
        await commandClass.execute(textList.slice(1));
    }

    const userInfo = JSON.parse(getCookie("__userInfo"));
    let username, ip;

    if (userInfo == null) {
        username = "Visitor";
        ip = "127.0.0.1";
    } else {
        username = userInfo["username"];
        ip = userInfo["IP"]
    }
    await addNewTerminal(username, ip);
}