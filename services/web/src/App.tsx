import ChatMessage from './components/ChatMessage';
import BettingModal from './components/BettingModal';
import Game from './Game';
import './App.css';
import React, { useState, useEffect } from 'react';
import RewardsButton from './components/RewadsButton';

export default function App() : JSX.Element {
  
  const [username, setUsername] : [string, React.Dispatch<React.SetStateAction<string>>] = useState("");

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => setUsername(e.target.value)

  const login = async function (event: React.FormEvent<HTMLFormElement>) : Promise<void> {
    event.preventDefault();
    const response = await fetch('http://localhost:8090/newUser', {
      method: "POST", 
      headers: {
        'Accept': 'application/json'
      },
      mode: 'cors',
      body: JSON.stringify({ username: username })
    });
    const json = await response.json();
    console.log(json)
    if (json.error === "") 
      setUsername(username)
  }

  return (
    <div className=" pl-2 pt-2 pr-2" id="container" style={{width: '100%', height: '100vh'}}>
      {(username === "" ? 
        <div className='card' style={{width: '50%', margin: '0 auto'}}>
          <h5 className="card-title">Create a Username</h5>
          <form onSubmit={login}>
            <input onChange={onChange}/>
            <button type='submit'>Submit</button>
          </form>
        </div> :
        <Game />
      )}
    </div>
  );
}
