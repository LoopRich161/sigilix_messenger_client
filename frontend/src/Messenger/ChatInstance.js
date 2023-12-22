import React, {useState} from 'react';
import { ListGroup, ButtonGroup, Button, Modal, Form } from 'react-bootstrap';
import {usePopup} from "../Popup";
import sigilixService from "./SigilixService";

function RenameChatModal({ show, handleClose, handleRename, currentName }) {
    const [newName, setNewName] = useState(currentName);

    const onSubmit = () => {
        handleRename(newName);
        handleClose();
    };

    return (
        <Modal show={show} onHide={handleClose}>
            <Modal.Header closeButton>
                <Modal.Title>Переименовать чат</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <Form onSubmit={e => { e.preventDefault(); }}>
                    <Form.Group className="mb-3" >
                        <Form.Label>New Chat Name</Form.Label>
                        <Form.Control
                            type="text"
                            value={newName}
                            onChange={e => setNewName(e.target.value)}
                            placeholder="Введите новое имя чата"
                        />
                    </Form.Group>
                </Form>
            </Modal.Body>
            <Modal.Footer>
                <Button variant="secondary" onClick={handleClose}>
                    Отмена
                </Button>
                <Button variant="primary" onClick={onSubmit}>
                    ОК
                </Button>
            </Modal.Footer>
        </Modal>
    );
}

function ChatInstance({ chat, onChatSelect, activeContextMenuChatId, setActiveContextMenuChatId}) {
    const { showPopup } = usePopup();
    const [contextMenuPosition, setContextMenuPosition] = useState({ x: 0, y: 0 });
    const [showRenameModal, setShowRenameModal] = useState(false);

    const thisChat = chat;

    const isContextMenuOpen = activeContextMenuChatId === chat.id;

    const handleRightClick = event => {
        event.preventDefault();
        setActiveContextMenuChatId(chat.id);
        setContextMenuPosition({ x: event.clientX, y: event.clientY });
    };

    const handleContextMenuSelect = (action) => {
        setActiveContextMenuChatId(null);
        console.log(`${action} chat: ${chat.id} ${thisChat.id}`);
        // Add your logic for deleting or renaming the chat here
        if (action === 'Delete') {
            sigilixService.deleteChat(chat.id).then(
                _ => {
                    showPopup("Чат удалён", "ОК", "success");
                }
            ).catch(
                error => {
                    console.error("Error deleting chat:", error);
                    showPopup("Не удалось удалить чат", "ОК", "danger");
                }
            )
        }
        if (action === 'Rename') {
            setShowRenameModal(true);
        }
    };

    const handleClick = () => {
        setActiveContextMenuChatId(null);
        onChatSelect(chat);
    };

    const handleRenameChat = newName => {
        console.log(`Rename chat from ${chat.title} to ${newName}`);
        sigilixService.renameChat(chat.id, newName).then(
            _ => {
                showPopup("Чат переименован", "ОК", "success");
            }
        ).catch(
            error => {
                console.error("Error renaming chat:", error);
                showPopup("Не удалось переименовать чат", "ОК", "danger");
            }
        )
    };

    return (
        <>
            <ListGroup.Item
                as="li"
                className="d-flex justify-content-between align-items-start"
                onClick={handleClick}
                onContextMenu={handleRightClick}
                style={{ cursor: 'pointer', marginTop: '5px', borderRadius: '10px' }}
            >
                <div className="ms-2 me-auto">
                    <div className="fw-bold">{chat.title}</div>
                    {chat.mbLastMessageText()}
                </div>
                {/*<Badge bg="primary" pill>*/}
                {/*    14*/}
                {/*</Badge>*/}
            </ListGroup.Item>

            {isContextMenuOpen && (
                <div
                    style={{
                        position: 'absolute',
                        top: `${contextMenuPosition.y}px`,
                        left: `${contextMenuPosition.x}px`,
                        zIndex: 1000
                    }}
                    className="bg-light p-2 border rounded"
                    onMouseLeave={() => setActiveContextMenuChatId(null)}
                >
                    <ButtonGroup vertical>
                        <Button size="sm" variant={'danger'} onClick={() => handleContextMenuSelect('Delete')}>Удалить чат</Button>
                        <Button size="sm" onClick={() => handleContextMenuSelect('Rename')}>Переименовать чат</Button>
                    </ButtonGroup>
                </div>
            )}

            <RenameChatModal
                show={showRenameModal}
                handleClose={() => setShowRenameModal(false)}
                handleRename={handleRenameChat}
                currentName={chat.title}
            />
        </>
    )
}

export default ChatInstance;
