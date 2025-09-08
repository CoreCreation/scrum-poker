import VoteButton from "./panel/voteButton";

export default function Panel({ options, sendVote, lastVote, disabled }) {
  return (
    <div class="vote-panel-options">
      {options.map((option) => (
        <VoteButton
          option={option}
          sendVote={sendVote}
          lastVote={lastVote}
          disabled={disabled}
        />
      ))}
    </div>
  );
}
