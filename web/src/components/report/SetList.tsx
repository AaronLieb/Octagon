import React, { useState, useEffect } from 'react';
import { Set } from '../../types';
import { ListSearch } from './ListSearch';
import './SetList.css';

interface SetListProps {
  tournament: string;
  onSelectSet: (set: Set) => void;
}

export function SetList({ tournament, onSelectSet }: SetListProps) {
  const [sets, setSets] = useState<Set[]>([]);
  const [readyToCallSets, setReadyToCallSets] = useState<Set[]>([]);
  const [loading, setLoading] = useState(true);
  const [includeCompleted, setIncludeCompleted] = useState(false);

  useEffect(() => {
    fetchSets();
  }, [tournament, includeCompleted]);

  useEffect(() => {
    fetchReadyToCall();
  }, [tournament]);

  const fetchSets = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/sets?tournament=${tournament}&includeCompleted=${includeCompleted}`);
      const data = await response.json();
      setSets(data.sets || []);
    } catch (error) {
      console.error('Failed to fetch sets:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchReadyToCall = async () => {
    try {
      const response = await fetch(`http://localhost:8080/api/sets/ready?tournament=${tournament}`);
      const data = await response.json();
      setReadyToCallSets(data.sets || []);
    } catch (error) {
      console.error('Failed to fetch ready sets:', error);
    }
  };

  const fuzzyFilter = (query: string, items: Set[], getLabel: (item: Set) => string) => {
    if (!query) return items;
    const lowerQuery = query.toLowerCase();
    return items.filter(item => 
      item.player1.name.toLowerCase().includes(lowerQuery) ||
      item.player2.name.toLowerCase().includes(lowerQuery)
    );
  };

  if (loading) return <div className="set-list">Loading...</div>;

  return (
    <div className="set-list-container">
      <div className="set-list-main">
        <div style={{ marginBottom: '10px' }}>
          <label>
            <input
              type="checkbox"
              checked={includeCompleted}
              onChange={e => setIncludeCompleted(e.target.checked)}
            />
            {' '}Include completed sets
          </label>
        </div>
        <ListSearch
          items={sets}
          getLabel={set => `${set.player1.name} vs ${set.player2.name} (${set.round})`}
          onSelect={onSelectSet}
          filterFn={fuzzyFilter}
          placeholder="Search players..."
        />
      </div>
      <div className="ready-to-call">
        <h3>Ready to Call</h3>
        {readyToCallSets.length > 0 ? (
          readyToCallSets.slice(0, 10).map(set => (
            <div key={set.id} className="ready-set" onClick={() => onSelectSet(set)}>
              <div className="ready-players">
                {set.player1.name} vs {set.player2.name}
              </div>
              <div className="ready-round">{set.round}</div>
            </div>
          ))
        ) : (
          <div style={{ color: '#666', fontSize: '14px' }}>No sets ready</div>
        )}
      </div>
    </div>
  );
}
