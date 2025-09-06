import VoteButton from "./panel/voteButton";

export default function Panel({ options, sendVote }) {
  return (
    <div>
      {options.map((option) => (
        <VoteButton option={option} sendVote={sendVote} />
      ))}
    </div>
  );
}
