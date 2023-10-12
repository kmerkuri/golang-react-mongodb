import React, { useState } from "react";

import Form from "react-bootstrap/Form";

import Button from "react-bootstrap/Button";


export default function Login() {

    const [form, setForm] = useState(
        {email: '', password: ''}
    )

    function validateForm() {

        return form.email.length > 0 && form.password.length > 0;

    }

    function handleSubmit(event) {

        event.preventDefault();
        fetch("http://localhost:8080/api/user/login", {
            method: "POST",
            body: JSON.stringify(form),
        }).then((response) => response.json())
            .then((data) => {
                console.log(data);
                // Handle data
            })
            .catch((err) => {
                console.log(err.message);
            });
    }

    return (

        <div className="Login">

            <Form onSubmit={handleSubmit}>

                <Form.Group size="lg" controlId="email">

                    <Form.Label>Email</Form.Label>

                    <Form.Control

                        autoFocus

                        type="email"

                        value={form.email}

                        onChange={(e) => setForm({...form, email: e.target.value})}

                    />

                </Form.Group>

                <Form.Group size="lg" controlId="password">

                    <Form.Label>Password</Form.Label>

                    <Form.Control

                        type="password"

                        value={form.password}

                        onChange={(e) => setForm({...form, password: e.target.value})}

                    />

                </Form.Group>

                <Button block size="lg" type="submit" disabled={!validateForm()}>

                    Login

                </Button>

            </Form>

        </div>

    );

}