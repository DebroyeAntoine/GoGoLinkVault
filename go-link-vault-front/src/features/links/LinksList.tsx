import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { fetchLinks } from './linksSlice';
import { RootState } from '../../app/store';

const LinksList = () => {
  const dispatch = useDispatch();
  const links = useSelector((state: RootState) => state.links.links);
  const status = useSelector((state: RootState) => state.links.status);
  const error = useSelector((state: RootState) => state.links.error);

  useEffect(() => {
    if (status === 'idle') {
      dispatch(fetchLinks());
    }
  }, [dispatch, status]);

  return (
    <div>
      <h1>Links</h1>
      {status === 'loading' ? (
        <p>Loading...</p>
      ) : status === 'failed' ? (
        <p>Error: {error}</p>
      ) : (
        <ul>
          {links.map((link: any) => (
            <li key={link.id}>
              <a href={link.url}>{link.title}</a>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default LinksList;

