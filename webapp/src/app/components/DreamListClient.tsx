import React from "react";

interface Props {
  dreams: Dream[];
}

export default function DreamClientList({ dreams }: Props) {
  return (
    <div>
      <h2>Current Dreams</h2>
      <ul>
        {dreams.map((dream) => (
          <li key={dream.id}>{dream.dream}</li>
        ))}
      </ul>
    </div>
  );
}
