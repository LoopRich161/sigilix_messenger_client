import React, {useEffect, useState} from 'react';
import {Container, Row, Col} from 'react-bootstrap';
import ChatsList from "./ChatView";
import Dialog from "./DialogComponent";
import {Navbar} from "react-bootstrap";
import {Nav} from "react-bootstrap";
import {NavDropdown} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import InputGroup from "react-bootstrap/InputGroup";

import sigilixService from "./SigilixService";
import Button from 'react-bootstrap/Button';
import ButtonGroup from 'react-bootstrap/ButtonGroup';
import {usePopup} from "../Popup";

function MainSplitView() {
    const [chats, setChats] = useState([]);
    const [currentChat, setCurrentChat] = useState(null);
    const { showPopup } = usePopup();

    useEffect(() => {
        const loadChats = async () => {
            try {
                const fetchedChats = await sigilixService.fetchChats();
                console.log("Fetched chats:", fetchedChats);
                setChats(fetchedChats);
            } catch (error) {
                console.error("Error loading chats:", error);
                // Handle error appropriately
            }
        };

        loadChats();
    }, [setChats]);

    const handleChatSelect = (chatId) => {
        setCurrentChat(chatId);
    };

    const addMessageToCurrentChat = (message) => {
        if (!currentChat) return;

        // Update current chat with new message
        currentChat.addMessage(message);

        // Update chats state
        setChats(chats.map(chat => chat.id === currentChat.id ? currentChat : chat));

        // Force update of current chat
        setCurrentChat(currentChat);

        // Hack to scroll to bottom of messages
        setTimeout(
            () => {
                const element = document.getElementById("message-container");
                element.scrollTop = element.scrollHeight;
            },
            200
        )
    };

    const requestChat = () => {
        const element = document.getElementById("request-chat-username");
        const username = element.value;

        if (!username) {
            showPopup("Введите имя пользователя", "ОК", "error");
            return;
        }
        sigilixService.TryRequestChat(username).then(
            _ => {
                element.value = "";
                showPopup("Вы запросили чат", "ОК", "success");
            }
        ).catch(
            error => {
                element.value = "";
                showPopup("Ошибка при запросе чата: " + error, "ОК", "error");
            }
        );
    }

    const copyUserId = () => {
        navigator.clipboard.writeText(sigilixService.userId).then(
            _ => {
                showPopup("ID скопирован в буфер обмена", "ОК", "success");
            }
        ).catch(
            error => {
                showPopup("Ошибка при копировании ID: " + error, "ОК", "error");
            }
        );
    }

    sigilixService.setCurrentOpenChat(currentChat?.id, addMessageToCurrentChat);
    sigilixService.setChatsCallback(setChats);


    return (
        // container, fixed height, 100% height, not scrollable
        <div style={{height: '100vh', overflow: 'hidden'}}>
            <Navbar expand="lg" className="bg-body-tertiary" color={'transparent'} sticky={'top'} style={{background: 'linear-gradient(to right, #0b95e850, #aed2d9a0)'}}>
                <Container>
                    <Navbar.Brand>Sigilix</Navbar.Brand>
                    <Navbar.Toggle aria-controls="basic-navbar-nav" />
                    <Navbar.Collapse id="basic-navbar-nav">
                        <Nav className="me-auto">
                            <Form inline onSubmit={e => { e.preventDefault(); }}>
                                <Row>
                                    <Col xs="auto">
                                        <Form.Control
                                            type="text"
                                            placeholder="Пользователь"
                                            className=" me-2"
                                            id={'request-chat-username'}
                                        />
                                    </Col>
                                    <Col xs="auto">
                                        <Button type="submit" onClick={requestChat}> Запросить диалог </Button>
                                    </Col>
                                    <Col xs="auto">
                                        <Button onClick={copyUserId}> Скопировать мой ID </Button>
                                    </Col>
                                </Row>

                            </Form>
                        </Nav>
                    </Navbar.Collapse>
                </Container>
            </Navbar>


            <Container fluid>
                <Row>
                    <Col md={3} style={{overflowY: 'auto', minHeight: '95vh', maxHeight: '95vh', background: 'linear-gradient(to right, #c899e7, #dfeedd)'}}>
                        <ChatsList chats={chats} onChatSelect={handleChatSelect}/>
                    </Col>
                    <Col md={9} className="d-flex flex-column" style={{maxHeight: '95vh', background: 'linear-gradient(to right, #dfeedd, #ecd4bf)'}}>
                        {<Dialog currentChat={currentChat} addMessage={addMessageToCurrentChat}/>}
                    </Col>
                </Row>
            </Container>
        </div>

    );
}

export default MainSplitView;