/*** 
 * <!-- Button trigger modal -->
<button type="button" className="btn btn-primary" data-toggle="modal" data-target="#exampleModalCenter">
  Launch demo modal
</button>
 */
import { useState } from 'react'

interface id {
    id: string
}

export default function BettingModal({ id }: id) : JSX.Element {

  const [amount, setAmount] : [number, React.Dispatch<React.SetStateAction<number>>] = useState(0);

  const handleChange = function (event: React.FormEvent<HTMLInputElement>) {
    let change: number = parseInt(event.currentTarget.value);
    if (!isNaN(change)) {
      setAmount(change)
    }  
  }

  const placeBet = async function () {

    if (amount === 0) return;

    let form = new URLSearchParams();
    form.set('address', '0xc0E5Af808cF0C15dfa145AF295A8F6B63DaE0258');
    form.set('amount', amount.toString());

    const response = await fetch("http://localhost:8000/bet", {
      method: 'POST', 
      mode: 'cors', 
      headers: {'Content-Type': 'application/x-www-form-urlencoded'},
      body: form
    });

    const obj = await response.json(); 
    console.log(obj)
    if (obj.error){
      alert("betting error")
    }
  }

  return (
    <div className="modal fade show" id={id}  tabIndex={-1} role="dialog" aria-labelledby={id} /*aria-hidden="true"*/>
    <div className="modal-dialog modal-dialog-centered" role="document">
        <div className="modal-content">
        <div className="modal-header">
          <h5 className="modal-title" id="exampleModalLongTitle">Modal title</h5>
          <button type="button" className="close" data-dismiss="modal" aria-label="Close">
          <span aria-hidden="true">&times;</span>
          </button>
        </div>
        <div className="modal-body">
          <input type="text" 
            className="form-control" 
            name="amount" 
            value={amount} 
            onChange={handleChange} 
            id="amount" 
            aria-label="Sizing example input" 
            aria-describedby="inputGroup-sizing-sm" 
          />
        </div>
        <div className="modal-footer">
          <button type="button" className="btn btn-secondary" onClick={placeBet} data-dismiss="modal">Place Bet</button>
          <button type="button" className="btn btn-primary" data-dismiss="modal">Cancel</button>
        </div>
        </div>
    </div>
    </div>
  );
}


