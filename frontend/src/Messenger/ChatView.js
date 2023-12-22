import React, {useState} from 'react';
import ChatInstance from './ChatInstance';
import ListGroup from 'react-bootstrap/ListGroup';

function ChatsList({ chats, onChatSelect}) {
    const [activeContextMenuChatId, setActiveContextMenuChatId] = useState(null);

    return (
        <div className="chats-list">
            <ListGroup as="ol">
                {chats.map(chat => <ChatInstance key={chat.id} chat={chat} onChatSelect={onChatSelect} activeContextMenuChatId={activeContextMenuChatId} setActiveContextMenuChatId={setActiveContextMenuChatId} />)}
            </ListGroup>
        </div>
    );
}

export default ChatsList;