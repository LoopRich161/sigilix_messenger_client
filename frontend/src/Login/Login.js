import Form from 'react-bootstrap/Form';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import {usePopup} from '../Popup';
import '../GenKeys/GenKeys.css';
import React from 'react';
import sigilixService from "../Messenger/SigilixService";

function Login({ switchView }) {
    const { showPopup } = usePopup();

    const handleSignup = () => {
        const loginButton = document.getElementById('login');
        loginButton.disabled = true;
        loginButton.innerHTML = '<span class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span> Загрузка...';

        const reversebutton = () => {
            loginButton.innerHTML = 'Разблокировать';
            loginButton.disabled = false;
        }

        sigilixService.login(document.getElementById('login-password').value).then(
            response => {
                reversebutton();
                if (!response) {
                    showPopup('Неверный пароль', 'Закрыть', 'danger');
                } else {
                    switchView('messenger');
                }
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
                    <Form.Label>Введите пароль</Form.Label>
                    <Form.Control type="password" placeholder="Пароль" id={'login-password'} />
                    <Button variant="primary" type="submit" className="mt-3" style={{float: "right"}} onClick={handleSignup} id={"login"} >
                        Разблокировать
                    </Button>
                </Form>
            </Container>
        </div>
    );
}

export default Login;