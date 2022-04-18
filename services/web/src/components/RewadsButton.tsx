export default function RewardsButton() : JSX.Element {
  const collectRewards = async function () {
    let form = new URLSearchParams();
    form.set('address', '0xc0E5Af808cF0C15dfa145AF295A8F6B63DaE0258');
    const response = await fetch("http://localhost:8000/collect", {
      method: 'POST', 
      mode: 'cors', 
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
      },
      body: form
    });
    const obj = await response.json();  
    if (obj.error){
      alert("you have no rewards")
    }
  }
  return (
    <button onClick={collectRewards}>
      Collect Rewards
    </button>
  );
}


