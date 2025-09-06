export default function VoteButton({ option, sendVote }) {
  return <button onClick={() => sendVote(option)}>{option}</button>;
}
