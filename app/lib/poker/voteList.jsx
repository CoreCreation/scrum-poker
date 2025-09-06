export default function VoteList({ data }) {
  console.log(data);
  if (!data || !data.length) return <div>Waiting for Voters!</div>;
  const votes = data.map((i) => i.vote).filter((i) => i !== -1);
  const average = votes.reduce((p, c) => p + c, 0) / votes.length;
  return (
    <table>
      <thead>
        <tr>
          <th scope="col">User</th>
          <th scope="col">Vote</th>
        </tr>
      </thead>
      <tbody>{data.map(Row)}</tbody>
      <tfoot>
        <tr>
          <th scope="row">Average</th>
          <td>{!Number.isNaN(average) ? average : "n/a"}</td>
        </tr>
      </tfoot>
    </table>
  );
}

function Row({ name, vote }) {
  return (
    <tr>
      <th scope="row">{name}</th>
      <td>{vote === -1 ? "No Vote" : vote}</td>
    </tr>
  );
}
