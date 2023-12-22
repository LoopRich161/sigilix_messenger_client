import React, { createContext, useState, useContext } from 'react';
import { Modal, Button } from 'react-bootstrap';
import sigilixService from "./Messenger/SigilixService";

const Popup = ({text, buttonText = 'Close', buttonVariant = 'primary', onClose}) => {
    buttonText = buttonText || 'Close';
    buttonVariant = buttonVariant || 'primary';
    return (
        <Modal show={true} onHide={onClose} centered style={{backdropFilter: 'blur(5px)'}}>
            {/*<Modal.Header closeButton>*/}
            {/*    <Modal.Title>Popup</Modal.Title>*/}
            {/*</Modal.Header>*/}
            <Modal.Body>
                <p style={{textAlign: 'center', fontSize: '20px', marginTop: '4%'}}>{text}</p>
            </Modal.Body>
            <Modal.Footer style={{borderColor: "transparent", justifyContent: 'center'}}>
                <Button variant={buttonVariant} onClick={onClose}>
                    {buttonText}
                </Button>
            </Modal.Footer>
        </Modal>
    );
};

const PopupContext = createContext();

export const usePopup = () => useContext(PopupContext);

export const PopupProvider = ({ children }) => {
    const [popup, setPopup] = useState({ isVisible: false, text: '', buttonText: 'Close', buttonVariant: 'primary' });

    const showPopup = (text, buttonText, buttonVariant) => setPopup({ isVisible: true, text, buttonText, buttonVariant });
    const hidePopup = () => setPopup({ isVisible: false, text: '', buttonText: 'Close', buttonVariant: 'primary' });

    sigilixService.setPopUpCallback(showPopup)

    return (
        <PopupContext.Provider value={{ showPopup, hidePopup }}>
            {children}
            {popup.isVisible && <Popup text={popup.text} onClose={hidePopup} buttonText={popup.buttonText} buttonVariant={popup.buttonVariant} />}
        </PopupContext.Provider>
    );
};


export default Popup;
