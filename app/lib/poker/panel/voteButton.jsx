export default function VoteButton({ option, sendVote, lastVote, disabled }) {
  return (
    <button
      class={
        lastVote === option && !disabled
          ? "outline secondary selected"
          : "outline"
      }
      onClick={() => sendVote(option)}
      disabled={disabled}
    >
      {option}
    </button>
  );
}
