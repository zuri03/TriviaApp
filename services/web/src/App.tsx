import Game from './Game';
import './App.css';
import React, { useState } from 'react';

export default function App() : JSX.Element {
  
  const [username, setUsername] : [string, React.Dispatch<React.SetStateAction<string>>] = useState("");
  const [password, setPassword] : [string, React.Dispatch<React.SetStateAction<string>>] = useState("");
  const [isLoggedIn, setIsLoggedIn] : [boolean, React.Dispatch<React.SetStateAction<boolean>>] = useState<boolean>(false)

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => 
    e.target.name === 'username' ? setUsername(e.target.value) : setPassword(e.target.value)

  const login = async function (event: React.FormEvent<HTMLFormElement>) : Promise<void> {
    console.log(username)
    console.log(password)
    event.preventDefault();
    const response = await fetch('http://localhost:8081/user', {
      method: "POST", 
      headers: {
        'Content-Type': 'application/json'
      },
      mode: 'cors',
      body: JSON.stringify({ username: username, password: password })
    });
    const json = await response.json();
    console.log(json)
    setIsLoggedIn(true)
    /*
    if (json.error === "") 
      setIsLoggedIn(true)
      */
  }

  return (
    <div className=" pl-2 pt-2 pr-2" id="container" style={{width: '100%', height: '100vh'}}>
      {(isLoggedIn ? 
        <Game />
        :
        <div className='card' style={{width: '50%', margin: '0 auto'}}>
          <h5 className="card-title">Create a Username</h5>
          <div className='card-body d-flex justify-content-center'>
            <div className='col-lg-5 text-center'>
              <form onSubmit={login}>
                <div className='py-2'>
                  <label>Username:</label>
                  <input onChange={onChange} name='username'/>
                  <label>Password:</label>
                  <input onChange={onChange} name='password'/>
                </div>
                <button type='submit' className='btn btn-primary'>Submit</button>
              </form>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
