import React, { useState } from 'react';

function UserForm() {
    // const [name, setName] = useState('');
    // const [email, setEmail] = useState('');
    // const [password, setPassword] = useState('');

    const [form, setForm] = useState(
        {name: '', email: '', password: ''}
    )
    console.log(form);
    const handleSubmit = (event) => {
        event.preventDefault();
        // fetch("http://localhost:8080/api/user/signup", {
        //     method: "POST",
        //     body: {
        //        form
        //     }
        // }).then(r => console.log(r.body))

        fetch("http://localhost:8080/api/user/signup", {
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

    };

    return (
        <div>
            <h2>Add User</h2>
            <form onSubmit={handleSubmit}>
                <label>
                    Name:
                    <input type="text" value={form.name} onChange={e => setForm({...form, name: e.target.value})} />
                </label>
                <label>
                    Email:
                    <input type="email" value={form.email} onChange={e => setForm({...form, email: e.target.value})} />
                </label>
                <label>
                    Password:
                    <input type="password" value={form.password} onChange={e => setForm({...form, password: e.target.value})} />
                </label>
                <button type="submit">Submit</button>
            </form>
        </div>
    );
}

export default UserForm;