import React from 'react';
import ReactDOM from 'react-dom/client';
import './index.css';
import UserList from './UserList';
import UserForm from './UserForm';
import Login from "./Login";
import { BrowserRouter, Routes, Route } from "react-router-dom";

export default function App() {
    return (
        <BrowserRouter>
            <Routes>


                    <Route path="/" element={<UserList />} />
                        <Route path="signin" element={<Login />} />
                        <Route path="signup" element={<UserForm />} />



            </Routes>
        </BrowserRouter>
    );
}

const root = ReactDOM.createRoot(document.getElementById('root'));
root.render(<App />);