import React, { useState, useEffect } from 'react';
import { Set } from '../../types';
import { ListSearch } from './ListSearch';
import API_URL, { getAuthHeaders } from '../../config';
import { ErrorMessage } from '../common/LoadingError';
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
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setLoading(true);
    fetchSets();
  }, [tournament, includeCompleted]);

  useEffect(() => {
    fetchReadyToCall();
  }, [tournament]);

  const fetchSets = async () => {
    setError(null);
    try {
      const response = await fetch(`${API_URL}/api/sets?tournament=${tournament}&includeCompleted=${includeCompleted}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) throw new Error('Failed to fetch sets');
      const data = await response.json();
      setSets(data.sets || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const fetchReadyToCall = async () => {
    try {
      const response = await fetch(`${API_URL}/api/sets/ready?tournament=${tournament}`, {
        headers: getAuthHeaders()
      });
      if (!response.ok) throw new Error('Failed to fetch ready sets');
      const data = await response.json();
      setReadyToCallSets(data.sets || []);
    } catch (err) {
      console.error('Failed to fetch ready sets:', err);
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
        {error && <ErrorMessage message={error} />}
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
