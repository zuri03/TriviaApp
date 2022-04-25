import ChatMessage from './components/ChatMessage';
import BettingModal from './components/BettingModal';
import './App.css';
import React, { useState, useEffect } from 'react';
import RewardsButton from './components/RewadsButton';

export default function App() : JSX.Element {
  
  const [question, setQuestion] : [string, React.Dispatch<React.SetStateAction<string>>] = useState("No Question Yet");
  const [playerAnswer, setPlayerAnswer] : [string, React.Dispatch<React.SetStateAction<string>>] = useState("");
  const [phase, setPhase] : [number, React.Dispatch<React.SetStateAction<number>>] = useState(0);
  
  useEffect(() => {
    startConnection();
  }, [])

  const startConnection = function () : void {
    let socket = new WebSocket("ws://127.0.0.1:8080/ws");
    console.log("Attempting Connection...");

    socket.onopen = () : void => {
      console.log("Successfully Connected");
      socket.send("Hi From the Client!")
    };
    
    socket.onmessage = (event: MessageEvent) : void => {
      let state: any = JSON.parse(event.data)
      setPhase(state.phase)
      setQuestion(state.question)
      console.log(`new phase => ${state.phase}`)
      if (state.phase === 4){
        submitAnswer();
      }
    }

    socket.onclose = (event: CloseEvent) : void => {
      console.log("Socket Closed Connection: ", event);
      socket.send("Client Closed!")
    };

    socket.onerror = (error: Event) : void => console.log("Socket Error: ", error);
  }

  const submitAnswer = async function () {
    const finalAnswer = localStorage.getItem("answer")
    localStorage.setItem('answer', '');
    setPlayerAnswer('')
    const response = await fetch("http://localhost:8080/Submit", {
      method: "POST", 
      headers: {
        'Accept': 'application/json'
      },
      mode: 'cors',
      body: JSON.stringify({ address: '0xc0E5Af808cF0C15dfa145AF295A8F6B63DaE0258', answer: finalAnswer })
    });

    const json = await response.json()
    console.log(json)
  }

  const storeAnswer = (event: React.FormEvent<HTMLFormElement>) : void => {
    event.preventDefault()
    localStorage.setItem("answer", playerAnswer)
  }
  
  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => setPlayerAnswer(e.target.value)

  const modalId = "betting-modal";

  return (
    <div className=" pl-2 pt-2 pr-2" id="container" style={{width: '100%', height: '100vh'}}>
      <BettingModal id={modalId} />
      <div className="row" style={{height: '90%', width: '100%'}}>
        <div className="col-9">
          <div className="card"> 
            {question}
          </div>
          <div className="card"> 
            Current Answer: {playerAnswer}
          </div>
        </div>
        <div className="col h-100">
          <div className="card" style={{height: '98%'}}>
            <div className="card-header">
              CHAT
            </div>
            <div className="card-body" id="chat-container" >
              <ChatMessage />
              <ChatMessage />
              <ChatMessage />
            </div>
            <div className="card-footer">
              <input type="text" className="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-sm" />
            </div>
          </div>
        </div>
      </div>
      <div className="d-flex flex-row justify-content-space-betweeen p-2" style={{backgroundColor: 'rgb(11, 32, 66)'}}>
        <div className="input-group mb-3">
          <form style={{width: '70%'}} onSubmit={storeAnswer}>
            <div className='d-flex flex-row justify-content-center'>
              <span className="input-group-text" id="inputGroup-sizing-sm">Answer</span>
              <input type="text" className="form-control" name="answer"  value={playerAnswer} onChange={onChange} aria-label="Sizing example input" aria-describedby="inputGroup-sizing-sm" />
              <button type="submit" className='m-1 btn btn-primary'>Submit</button>
            </div>
          </form>
        </div>
        <div className="d-flex flex-row justify-content-space-between">
          <button type="button"  className="m-1 btn btn-primary" data-toggle="modal" data-target={`#${modalId}`}>
            Place Bet
          </button>
          <RewardsButton />
        </div> 
      </div>
    </div>
  );
}
