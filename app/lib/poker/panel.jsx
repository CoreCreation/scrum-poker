import VoteButton from "./panel/voteButton";

export default function Panel({ options, sendVote }) {
  return (
    <div class="vote-panel-options">
      {options.map((option) => (
        <VoteButton option={option} sendVote={sendVote} />
      ))}
    </div>
  );
}
