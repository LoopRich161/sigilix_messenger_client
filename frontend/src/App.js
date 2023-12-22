import './App.css';
import GenKeys from "./GenKeys/Genkeys";
import Login from "./Login/Login";
import {useEffect, useState} from "react";
import {usePopup} from "./Popup";
import MainSplitView from "./Messenger/MainSplitView";
import sigilixService from "./Messenger/SigilixService";

function Loading() {
    return (
        <div className="gradient-background d-flex justify-content-center align-items-center" style={{ height: '100vh' }}>
            <div className="spinner-border text-light" role="status">
                <span className="sr-only"></span>
            </div>
        </div>
    );
}

function App() {
    const [view, setView] = useState('loading');

    // load state from sigilixService
    useEffect(() => {
        const loadState = async () => {
            try {
                const fetchedState = await sigilixService.getState();
                console.log("Fetched state:", fetchedState);
                setView(fetchedState);
            } catch (error) {
                console.error("Error loading state:", error);
                // Handle error appropriately
            }
        };

        loadState();
    }, [setView]);

    return (
        <div>
            {view === 'loading' && <Loading/>}
            {view === 'signup' && <GenKeys switchView={setView}/>}
            {view === 'login' && <Login switchView={setView}/>}
            {view === 'messenger' && <MainSplitView switchView={setView}/>}
        </div>
    );
}

export default App;
