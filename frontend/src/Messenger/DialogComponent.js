import React, {useEffect, useRef} from 'react';
import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import InputGroup from 'react-bootstrap/InputGroup';
import FormControl from "react-bootstrap/FormControl";
import Container from 'react-bootstrap/Container';

import sigilixService, {Message} from "./SigilixService";
import {usePopup} from "../Popup";

function MessageBubble({ message }) {
    const bubbleStyle = {
        maxWidth: '70%', // Set a max width for bubbles
        padding: '10px',
        margin: '10px 0',
        borderRadius: '15px',
        wordWrap: 'break-word' // Ensures text wraps within the bubble
    };

    const alignmentClass = message.sentByUs ? 'justify-content-end' : 'justify-content-start';
    const backgroundColorClass = message.sentByUs ? 'bg-primary text-white' : 'bg-light';

    return (
        <div className={`d-flex ${alignmentClass}`}>
            <div style={bubbleStyle} className={`shadow-sm ${backgroundColorClass}`}>
                {message.text}
            </div>
        </div>
    );
}

function Dialog({currentChat, addMessage }) {
    const inputFieldEnabled = currentChat?.isAccepted;
    let inputFieldPlaceholder = "";
    const { showPopup } = usePopup();

    if (!currentChat) {
        inputFieldPlaceholder = "Выберите чат, чтобы начать общение";
    } else if (!currentChat.isAccepted) {
        if (currentChat.isCreator) {
            inputFieldPlaceholder = "Ваш собеседник ещё не принял чат";
        } else {
            inputFieldPlaceholder = "Вы ещё не приняли чат";
        }
    } else {
        inputFieldPlaceholder = "Введите сообщение";
    }

    function sendMessage() {
        if (!currentChat || !currentChat.isAccepted) {
            return;
        }
        console.log("Sending message");
        const element = document.getElementById("send-text");
        const text = element.value;
        if (!text) {
            return;
        }
        sigilixService.sendMessage(currentChat.id, text).then(
            _ => console.log("Message sent"),
        )
        console.log(text);
        // const message = new Message(currentChat.getLastMessage()?.id + 1, currentChat.id, true, text);
        // addMessage(message);
        element.value = "";
    }

    const messagesEndRef = useRef(null);

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "auto" });
    };

    useEffect(() => {
        scrollToBottom();
    }, [currentChat?.messages]);

    const canSend = currentChat?.isAccepted;

    const acceptChat = () => {
        sigilixService.acceptChat(currentChat.id).then(
            _ => {
                showPopup("Вы приняли чат", "ОК", "success");
            }
        ).catch(
            error => {
                console.error("Error accepting chat:", error);
                showPopup("Не удалось принять чат", "ОК", "danger");
            }
        )
    }

    // show accept chat button if chat is not accepted

    const canAccept = currentChat && !currentChat.isAccepted && !currentChat.isCreator;


    return (
        <Container className="mt-auto">
            <Container style={{overflowY: 'auto', maxHeight: '90vh', padding: '10px', scrollbarColor: '#ced4da #f8f9fa'}} id={'message-container'}>
                {currentChat?.messages.map(message => <MessageBubble key={message.id} message={message} />)}
                {
                    canAccept &&
                    <Container className="d-flex justify-content-center align-items-center" style={{height: '100%'}}>
                        <Button variant="primary" onClick={acceptChat} size={'lg'}>Принять чат</Button>
                    </Container>
                }
                <div ref={messagesEndRef} /> {/* Invisible element at the end of messages */}

            </Container>
            <Container style={{marginBottom: '10px'}}>
                <Form onSubmit={e => { e.preventDefault(); }}>
                    <InputGroup className="mb-3">
                        <FormControl
                            as="textarea"
                            placeholder={inputFieldPlaceholder}
                            disabled={!inputFieldEnabled}
                            // handle ctrl+enter
                            onKeyDown={(e) => {
                                if (e.key === 'Enter' && e.ctrlKey) {
                                    sendMessage();
                                }
                            }}
                            id="send-text"
                        />
                        <Button onClick={sendMessage} disabled={!canSend} variant={!canSend ? "outline-secondary" : "primary"} >
                            Отправить
                        </Button>
                    </InputGroup>
                </Form>
            </Container>
        </Container>
    );
}

export default Dialog;