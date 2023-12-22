import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import {usePopup} from '../Popup';
import './GenKeys.css';
import React from 'react';
import sigilixService from "../Messenger/SigilixService";

function GenKeys({ switchView }) {
    const { showPopup } = usePopup();

    const handleSignup = () => {
        const genKeysButton = document.getElementById('gen-keys');

        if (document.getElementById('signup-password').value !== document.getElementById('signup-password-confirm').value) {
            showPopup('Пароли не совпадают', 'Закрыть', 'danger');
            return;
        }

        // check length of password
        if (document.getElementById('signup-password').value.length === 0) {
            showPopup('Пароль не может быть пустым', 'Закрыть', 'danger');
            return;
        }

        // loading animation from bootstrap
        genKeysButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Генерация...';
        genKeysButton.disabled = true;

        const reversebutton = () => {
            genKeysButton.innerHTML = 'Сгенерировать ключи';
            genKeysButton.disabled = false;
        }

        sigilixService.signUp(document.getElementById('signup-username').value, document.getElementById('signup-password').value).then(
            _ => {
                reversebutton();
                switchView('login');
            }
        ).catch(
            error => {
                reversebutton();
                showPopup(error, 'Закрыть', 'danger');
            }
        );
    }

    return (
        <div className="gradient-background d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
            <Container className="p-md-3" style={{background: 'rgba(255, 255, 255, 0.6)', borderRadius: '10px'}}>
                <Form onSubmit={e => { e.preventDefault(); }}>
                    <Form.Label>Username (Оставьте пустым, если не хотите, чтобы вас можно было найти по нему)</Form.Label>
                    <Form.Control type="text" placeholder="Username" id={'signup-username'} />
                    <Form.Label>Придумайте пароль</Form.Label>
                    <Form.Control type="password" placeholder="Пароль" id={'signup-password'} />
                    <Form.Label  style={{marginTop: '5px'}} >Повтор пароля</Form.Label>
                    <Form.Control type="password" placeholder="Пароль" id={'signup-password-confirm'}/>
                    <Button variant="primary" type="submit" className="mt-3" style={{float: "right"}} onClick={handleSignup} id={"gen-keys"} >
                        Сгенерировать ключи
                    </Button>
                </Form>
            </Container>
        </div>
    );
}

export default GenKeys;