function GetChat(arg1) {
    return window['go']['main']['App']['GetChat'](arg1);
}

function GetChatMessages(arg1) {
    return window['go']['main']['App']['GetChatMessages'](arg1);
}

function GetChats() {
    return window['go']['main']['App']['GetChats']();
}

function GetState() {
    return window['go']['main']['App']['GetState']();
}

function GetUserId() {
    return window['go']['main']['App']['GetUserId']();
}

function GetUsername() {
    return window['go']['main']['App']['GetUsername']();
}

function Greet(arg1) {
    return window['go']['main']['App']['Greet'](arg1);
}

function InitChatFromInitializer(arg1) {
    return window['go']['main']['App']['InitChatFromInitializer'](arg1);
}

function InitChatFromReceiver(arg1) {
    return window['go']['main']['App']['InitChatFromReceiver'](arg1);
}

function IsSignedUp() {
    return window['go']['main']['App']['IsSignedUp']();
}

function IsUnlocked() {
    return window['go']['main']['App']['IsUnlocked']();
}

function PullNotificationsAndUpdateData() {
    return window['go']['main']['App']['PullNotificationsAndUpdateData']();
}

function SearchByUsername(arg1) {
    return window['go']['main']['App']['SearchByUsername'](arg1);
}

function SendMessage(arg1, arg2) {
    return window['go']['main']['App']['SendMessage'](arg1, arg2);
}

function SetUsernameConfig(arg1, arg2) {
    return window['go']['main']['App']['SetUsernameConfig'](arg1, arg2);
}

function SignUp(arg1) {
    return window['go']['main']['App']['SignUp'](arg1);
}

function Unlock(arg1) {
    return window['go']['main']['App']['Unlock'](arg1);
}

function TryRequestChat(arg1) {
    return window['go']['main']['App']['TryRequestChat'](arg1);
}

function DeleteChat(arg1) {
    return window['go']['main']['App']['DeleteChat'](arg1);
}


function RenameChat(arg1, arg2) {
    return window['go']['main']['App']['RenameChat'](arg1, arg2);
}

class Chat {
    constructor(id, title, isCreator, isAccepted, messages, otherUserId) {
        this.id = id;
        this.title = title;
        this.isCreator = isCreator;
        this.isAccepted = isAccepted;
        this.messages = messages;
        this.otherUserId = otherUserId;
    }

    addMessage(message) {
        // check type of message
        if (message.ismsg === undefined) {
            throw new Error('Invalid message type');
        }
        this.messages.push(message);
    }

    getMessages() {
        return this.messages;
    }

    getLastMessage() {
        if (this.messages.length === 0) {
            return null;
        }
        return this.messages[this.messages.length - 1];
    }

    mbLastMessageText(trimLength = 20) {
        const lastMessage = this.getLastMessage();
        if (!lastMessage) {
            return '';
        }
        const text = lastMessage.text;
        if (text.length <= trimLength) {
            return text;
        }
        return text.substring(0, trimLength) + '...';
    }
}

class Message {
    constructor(id, chatId, sentByUs, text) {
        this.id = id;
        this.sentByUs = sentByUs;
        this.text = text;
    }

    ismsg() {
    }
}

class SigilixService {

    constructor() {
        this.chatStorage = new Map();
        this.loggedIn = false;
        this.userId = 0;
        this.username = '';

        this.currentOpenChatId = null;
        this.currentOpenChatAddMessageCallback = null; // shall be called if new message is received

        this.chatsCallback = null; // shall be called if chats are updated

        this.popUpCallback = null; // shall be called if popUp is updated

        this.notificationPullerRunning = false;

        IsUnlocked().then(
            a => {
                this.loggedIn = a;
                if (a) {
                    // make notificationPuller run in the background
                    GetUserId().then(
                        id => {
                            this.userId = id;
                        }
                    )
                    this.notificationPuller();
                }
            }
        )


    }

    setPopUpCallback(callback) {
        this.popUpCallback = callback;
    }

    showErrorPopUp(error) {
        this.popUpCallback?.(error, 'Закрыть', 'danger');
    }

    setCurrentOpenChat(chatId, addMessageCallback) {
        this.currentOpenChatId = chatId;
        this.currentOpenChatAddMessageCallback = addMessageCallback;
    }

    setChatsCallback(callback) {
        this.chatsCallback = callback;
    }

//     async fetchChats() {
//         // Fetch the list of chats from the API
//         // TODO: call API to fetch chats
//
//         // create mock chats
//         const newChats = [
//             new Chat(1, 'Chat 1', true, true, [
//                 new Message(1, 1, true, 'Hello'),
//                 new Message(2, 1, false, 'Hi'),
//                 new Message(3, 1, true, 'How are you?'),
//                 new Message(4, 1, false, 'Good, thanks'),
//                 new Message(5, 1, true, 'What are you doing?'),
//                 new Message(6, 1, false, 'Nothing special, wby?'),
//                 new Message(7, 1, true, 'Same'),
//                 new Message(8, 1, false, 'Ok. How\'s your family?'),
//                 new Message(9, 1, true, 'They are fine'),
//                 new Message(10, 1, false, 'Good. Anything new on the job? New framework?'),
//                 new Message(11, 1, true,
// `Yeah, we're using React now
// It's pretty cool
// I'm learning it now.
// So far I've learned about components, props, state, hooks, and some other stuff, but I'm still not sure how to use it all together. I think I'll need to build something to get a better understanding of how it all works together.
// `
//                 ),
//                 new Message(12, 1, false, 'Sounds cool'),
//                 new Message(13, 1, true, 'Yeah, it is'),
//                 new Message(14, 1, false, 'Good luck with that'),
//                 new Message(1411, 1, false, 'I\'ve started playing a new game recently'),
//                 new Message(15, 1, true, 'What game?'),
//                 new Message(16, 1, false, 'It\'s called Cyberpunk 2077'),
//                 new Message(17, 1, false, 'It\'s pretty cool, but it\'s a bit buggy'),
//                 new Message(18, 1, true, 'I\'ve heard about it'),
//                 new Message(19, 1, true, 'I\'ve heard it\'s buggy'),
//                 new Message(20, 1, false, 'Yeah, it is'),
//                 new Message(21, 1, false, 'But it\'s still fun'),
//                 new Message(22, 1, true, 'I\'ll check it out'),
//                 new Message(23, 1, false, 'Cool'),
//                 new Message(24, 1, false, 'Well, I gotta go now'),
//                 new Message(25, 1, true, 'Ok, bye'),
//                 new Message(26, 1, false, 'Bye'),
//             ], 2),
//             new Chat(2, 'Chat 2', false, true, [
//                 new Message(1, 2, false, 'Hello'),
//                 new Message(2, 2, true, 'Hi'),
//                 new Message(3, 2, false, 'How are you?'),
//                 new Message(4, 2, true, 'Good, thanks'),
//             ], 3),
//             new Chat(3, 'Chat 3', false, false, [], 4),
//             new Chat(4, 'Chat 4', true, false, [], 4),
//         ]
//
//         for (let i = 5; i < 105; i++) {
//             newChats.push(new Chat(i, `Chat ${i}`, true, false, [], 4));
//         }
//
//         for (const chat of newChats) {
//             this.chatStorage.set(chat.id, chat)
//         }
//         return this.arrayOfChats();
//     }

    async fetchChats() {

        const chats = await GetChats();
        console.log("Fetched chats:", chats);
        for (const chat of chats) {
            this.chatStorage.set(chat.chat_id, this.dataChatToChat(chat));
        }
        return this.arrayOfChats();
    }

    dataMessageToMessage(dataMessage) {
        return new Message(dataMessage.message_id, dataMessage.chat_id, dataMessage.sender_id === this.userId, dataMessage.content);
    }

    dataChatToChat(dataChat) {

        let messages = [];
        if (dataChat.messages) {
            messages = dataChat.messages.map(message => this.dataMessageToMessage(message));
        }

        return new Chat(dataChat.chat_id, dataChat.title, dataChat.am_i_initiator, dataChat.accepted, messages, dataChat.other_user_id);
    }

    mbDataChatToExistingChat(dataChat) {
        const chat = this.chatStorage.get(dataChat.chat_id);
        const chatFromData = this.dataChatToChat(dataChat);
        if (!chat) {
            return chatFromData;
        }

        // update fields
        chat.title = chatFromData.title;
        chat.isCreator = chatFromData.isCreator;
        chat.isAccepted = chatFromData.isAccepted;
        chat.otherUserId = chatFromData.otherUserId;
        chat.messages = chatFromData.messages;

        return chat;
    }

    async sendMessage(chatId, messageText) {
        if (!this.chatStorage.has(chatId)) {
            throw new Error('Chat not found');
        }

        SendMessage(chatId, messageText).then(
            msg => {
                const message = this.dataMessageToMessage(msg);
                if (this.currentOpenChatId === chatId) {
                    this.currentOpenChatAddMessageCallback?.(message);
                } else {
                    this.chatStorage.get(chatId).addMessage(message);
                }
            }
        ).catch(
            error => this.showErrorPopUp(error)
        );
    }

    async acceptChat(chatId) {
        if (!this.chatStorage.has(chatId)) {
            // throw new Error('Chat not found');
            this.showErrorPopUp(new Error('Chat not found'));
            return;
        }
        const chat = this.chatStorage.get(chatId);
        if (chat.isAccepted) {
            // throw new Error('Chat already accepted');
            this.showErrorPopUp(new Error('Chat already accepted'));
            return;
        }
        if (chat.isCreator) {
            // throw new Error('Only receiver can accept chat');
            this.showErrorPopUp(new Error('Only receiver can accept chat'));
            return;
        }
        try {
            await InitChatFromReceiver(chatId);
            chat.isAccepted = true;
            this.chatsCallback?.(this.arrayOfChats());

        } catch (error) {
            this.showErrorPopUp(error);
        }
    }

    async createChat(userId) {

        // create mock chat
        // const newChat = new Chat(5, title, true, false, []);
        // this.chatStorage.set(newChat.id, newChat);

        const chat = await InitChatFromInitializer(userId);
        const newChat = this.dataChatToChat(chat);
        this.chatStorage.set(newChat.id, newChat);
        this.chatsCallback?.(this.arrayOfChats());
    }

    async TryRequestChat(usernameOrId) {
        const chat = await TryRequestChat(usernameOrId);
        const newChat = this.dataChatToChat(chat);
        this.chatStorage.set(newChat.id, newChat);
        this.chatsCallback?.(this.arrayOfChats());
    }

    async signUp(username, password) {
        await SignUp(password)
        await Unlock(password)
        if (username) {
            await SetUsernameConfig(username, true);
        }
        this.userId = await GetUserId();
        this.username = await GetUsername();
        this.loggedIn = true;
        return true;
    }

    async login(password) {
        await Unlock(password)
        this.userId = await GetUserId();
        this.username = await GetUsername();
        this.loggedIn = true;
        this.notificationPuller();
        return true;
    }


    lastChatId() {
        if (this.chatStorage.length === 0) {
            return -1;
        }
        var maxId = 0;
        for (const chat of this.chatStorage) {
            if (chat.id > maxId) {
                maxId = chat.id;
            }
        }
        return maxId;
    }

    arrayOfChats() {
        return Array.from(this.chatStorage.values());
    }

    async getState() {
        return await GetState();
    }

    async notificationPuller() {
        // pool notifications every 100ms
        if (this.notificationPullerRunning) {
            return
        }
        this.notificationPullerRunning = true;

        while (true) {
            console.log("Pulling notifications")
            await new Promise(r => setTimeout(r, 100));
            try {
                await this.pullNotificationsInner();
            } catch (e) {
                console.error("Error pulling notifications:", e);
            }
        }
    }

    async pullNotificationsInner() {
        const updates = await PullNotificationsAndUpdateData();

        for (const update of updates) {
            console.log("Update:", update);
            const notification = update.notification;
            if (update.type === "new_incoming_chat") {
                const chat = this.mbDataChatToExistingChat(notification.chat);

                this.chatStorage.set(chat.id, chat);
                this.chatsCallback?.(this.arrayOfChats());

            } else if (update.type === "new_message") {
                const chatId = notification.chat_id;
                const message = this.dataMessageToMessage(notification.message);
                const chat = this.chatStorage.get(chatId);
                if (this.currentOpenChatId === chat.id) {
                    this.currentOpenChatAddMessageCallback?.(message);
                } else {
                    chat.addMessage(message);
                }

            } else if (update.type === "chat_accepted") {
                const chat = this.mbDataChatToExistingChat(notification.chat);

                chat.isAccepted = true;
                this.chatsCallback?.(this.arrayOfChats());

            } else {
                console.error("Unknown update type:", update.type);
            }
        }
    }

    async deleteChat(chatId) {
        await DeleteChat(chatId);
        this.chatStorage.delete(chatId);
        this.chatsCallback?.(this.arrayOfChats());
    }

    async renameChat(chatId, title) {
        await RenameChat(chatId, title);
        const chat = this.chatStorage.get(chatId);
        chat.title = title;
        this.chatsCallback?.(this.arrayOfChats());
    }

    debug_AcceptChat(chatId) {
        // only for testing purposes
        if (!this.chatStorage.has(chatId)) {
            throw new Error('Chat not found');
        }
        const chat = this.chatStorage.get(chatId);
        if (chat.isAccepted) {
            throw new Error('Chat already accepted');
        }
        chat.isAccepted = true;

        this.chatsCallback?.(this.arrayOfChats());
    }

    debug_CreateChat(title, isCreator, isAccepted, messages, otherUserId) {
        // only for testing purposes
        const newChatId = this.lastChatId() + 1;
        const newChat = new Chat(newChatId, title, isCreator, isAccepted, messages, otherUserId);
        this.chatStorage.set(newChat.id, newChat);

        this.chatsCallback?.(this.arrayOfChats());
    }

    debug_ReceiveMessage(chatId, messageText) {
        // only for testing purposes
        if (!this.chatStorage.has(chatId)) {
            throw new Error('Chat not found');
        }
        const chat = this.chatStorage.get(chatId);

        const message = new Message(chat.getLastMessage()?.id + 1, chat.id, false, messageText);

        if (this.currentOpenChatId === chatId) {
            this.currentOpenChatAddMessageCallback?.(message);
        }
    }


    // Other API methods...
}

const sigilixService = new SigilixService();
export default sigilixService;
export {Chat, Message};

// export sigilixService to the global scope for debugging purposes
window.sigilixService = sigilixService;