export default function VoteButton({ option, sendVote, lastVote }) {
  return (
    <button
      class={lastVote === option ? "outline secondary selected" : "outline"}
      onClick={() => sendVote(option)}
    >
      {option}
    </button>
  );
}
