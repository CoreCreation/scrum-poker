export default function VoteButton({ option, sendVote }) {
  return (
    <button class="outline" onClick={() => sendVote(option)}>
      {option}
    </button>
  );
}
