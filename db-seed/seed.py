from sqlalchemy import create_engine, Column, String, ForeignKey, Numeric, DECIMAL
from sqlalchemy.pool import QueuePool
from sqlalchemy.orm import declarative_base, sessionmaker
import uuid
from faker import Faker
import random
from tqdm import tqdm
import argparse
import json

fake = Faker()
Base = declarative_base()


# pool_size = 20
# max_overflow = 25
# No need to define pool_size and max_overflow here, they will be taken from the config

class User(Base):
    __tablename__ = 'users'
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String)


class Event(Base):
    __tablename__ = 'events'
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    name = Column(String)
    type = Column(String)
    status = Column(String)


class Purchase(Base):
    __tablename__ = 'purchases'
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    user_id = Column(String, ForeignKey('users.id'))
    event_id = Column(String, ForeignKey('events.id'))
    status = Column(String)


class Payment(Base):
    __tablename__ = 'payments'
    id = Column(String, primary_key=True, default=lambda: str(uuid.uuid4()))
    purchase_id = Column(String, ForeignKey('purchases.id'))
    amount = Column(DECIMAL(10, 2))
    status = Column(String)


def create_database_engine(db_config, ssl_config, secure):
    if secure:
        ssl_args = f"?sslmode=verify-full&sslrootcert={ssl_config['rootcert']}&sslcert={ssl_config['certfile']}&sslkey={ssl_config['keyfile']}"
    else:
        ssl_args = "?sslmode=disable"

    connection_string = f"cockroachdb://{db_config['user']}:{db_config['password']}@" \
                        f"{db_config['host']}:{db_config['port']}/{db_config['database']}{ssl_args}"

    engine = create_engine(
        connection_string,
        poolclass=QueuePool,
        pool_size=db_config['pool_size'],
        max_overflow=db_config['pool_max_conns'] - db_config['pool_size'],
        pool_pre_ping=True
    )
    return engine


def create_fake_data(session, num_users, num_purchases, num_events, num_cancellations, num_payments, batch_size=100):
    user_ids = []
    event_ids = []
    purchase_ids = []

    for _ in tqdm(range(num_users), desc='Creating Users'):
        user = User(name=fake.name())
        user_ids.append(user)
    session.add_all(user_ids)
    session.commit()

    for _ in tqdm(range(num_events), desc='Creating Events'):
        event = Event(name=fake.sentence(nb_words=4), type=random.choice(
            ['concert', 'opera', 'theater', 'football', 'soccer', 'baseball', 'basketball', 'conference']),
                      status='scheduled')
        event_ids.append(event)
    session.add_all(event_ids)
    session.commit()

    for _ in tqdm(range(num_purchases), desc='Creating Purchases'):
        purchase = Purchase(user_id=random.choice(user_ids).id, event_id=random.choice(event_ids).id,
                            status='confirmed')
        purchase_ids.append(purchase)
    session.add_all(purchase_ids)
    session.commit()

    for _ in tqdm(range(num_cancellations), desc='Creating Cancellations'):
        purchase_id = random.choice(purchase_ids).id
        purchase = session.query(Purchase).filter_by(id=purchase_id).first()
        purchase.status = 'cancelled'
    session.commit()

    payment_batch = []
    for _ in tqdm(range(num_payments), desc='Creating Payments'):
        purchase_id = random.choice(purchase_ids).id
        payment = Payment(purchase_id=purchase_id, amount=random.uniform(20, 200), status='successful')
        payment_batch.append(payment)
        if len(payment_batch) >= batch_size:
            session.add_all(payment_batch)
            session.commit()
            payment_batch = []
    session.add_all(payment_batch)
    session.commit()


def main(args):
    with open(args.config) as f:
        config = json.load(f)
    print(config)
    db_config = config['cockroach']['database']
    ssl_config = config['cockroach']['ssl']
    print(db_config)
    engine = create_database_engine(db_config, ssl_config, args.secure)
    Base.metadata.create_all(engine)

    Session = sessionmaker(bind=engine)
    session = Session()

    create_fake_data(session, args.num_users, args.num_purchases, args.num_events, args.num_cancellations,
                     args.num_payments)

    session.close()
    print("Data generation complete!")


if __name__ == "__main__":
    # python db-seed/seed.py (--secure or null) --config config/config.json --num_users 1 --num_purchases 1 --num_events 1 --num_cancellations 1 --num_payments 1
    parser = argparse.ArgumentParser(description='Generate fake data for ticket purchasing service.')
    parser.add_argument('--secure', action='store_true', help='Use secure database connection')
    parser.add_argument('--config', type=str, default='config/config.json', help='Path to configuration file')
    parser.add_argument('--num_users', type=int, default=1000, help='Number of users to generate')
    parser.add_argument('--num_purchases', type=int, default=5000, help='Number of purchases to generate')
    parser.add_argument('--num_events', type=int, default=5000, help='Number of purchases to generate')
    parser.add_argument('--num_cancellations', type=int, default=1000, help='Number of cancellations to generate')
    parser.add_argument('--num_payments', type=int, default=5000, help='Number of payments to generate')

    args = parser.parse_args()
    print(args)
    main(args)
