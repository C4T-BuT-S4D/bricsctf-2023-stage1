
CREATE OR REPLACE FUNCTION aaa(b varbit) RETURNS bytea 
    LANGUAGE plpgsql AS $$
DECLARE
    uwu bytea := ''::bytea;
BEGIN
    FOR i IN REVERSE length(b)..1 BY 32 LOOP
        uwu := decode(lpad(to_hex(substring(b FROM i-31 FOR 32)::int), 8, '0'), 'hex') || uwu;
    END LOOP;
    RETURN uwu;
END;
$$;

CREATE OR REPLACE FUNCTION aab(b bytea) RETURNS varbit 
    LANGUAGE plpgsql AS $$
BEGIN
    RETURN right(b::text, -1)::varbit;
END;
$$;

CREATE OR REPLACE FUNCTION aba(a bytea, b bytea) RETURNS bytea
    LANGUAGE plpgsql AS $$
BEGIN
    RETURN aaa(aab(a)::bit(128) # aab(b)::bit(128));
END;
$$;

CREATE OR REPLACE FUNCTION baa(b bytea) RETURNS bytea 
    LANGUAGE plpgsql AS $$
DECLARE
    owo bytea := '\x637c777bf26b6fc53001672bfed7ab76ca82c97dfa5947f0add4a2af9ca472c0b7fd9326363ff7cc34a5e5f171d8311504c723c31896059a071280e2eb27b27509832c1a1b6e5aa0523bd6b329e32f8453d100ed20fcb15b6acbbe394a4c58cfd0efaafb434d338545f9027f503c9fa851a3408f929d38f5bcb6da2110fff3d2cd0c13ec5f974417c4a77e3d645d197360814fdc222a908846eeb814de5e0bdbe0323a0a4906245cc2d3ac629195e479e7c8376d8dd54ea96c56f4ea657aae08ba78252e1ca6b4c6e8dd741f4bbd8b8a703eb5664803f60e613557b986c11d9ee1f8981169d98e949b1e87e9ce5528df8ca1890dbfe6426841992d0fb054bb16'::bytea;
    uwu bytea := b;
BEGIN
    FOR i IN 0..length(b)-1 LOOP
        uwu := set_byte(uwu, i, get_byte(owo, get_byte(b, i)));
    END LOOP;
    RETURN uwu;
END;
$$;

-- verified, correct!
CREATE OR REPLACE FUNCTION bab(nya bytea) RETURNS bytea[]
    LANGUAGE plpgsql AS $$
DECLARE
    pnq int;
    qop int;
    qwq bit(32)[];
    unu bit(32)[];

    eoe bit(32)[] := array[1, 2, 4, 8, 16, 32, 64, 128, 27, 54]::bit(32)[];
BEGIN
    IF length(nya) = 16 THEN
        pnq := 4; qop := 11;
    ELSIF length(nya) = 24 THEN
        pnq := 6; qop := 13;
    ELSIF length(nya) = 32 THEN
        pnq := 8; qop := 15;
    ELSE
        RAISE EXCEPTION 'die nya';
    END IF;

    FOR i IN 0..pnq-1 LOOP
        qwq[i] := aab(substring(nya FROM i*4+1 FOR 4))::bit(32);
    END LOOP;


    FOR i IN 0..(4*qop-1) LOOP
        IF i < pnq THEN
            unu[i] := qwq[i];
        ELSE
            CASE
                WHEN i % pnq = 0 THEN
                    DECLARE
                        wee bit(32);
                    BEGIN
                        wee := (unu[i-1] << 8) | (unu[i-1] >> 24);
                        unu[i] := unu[i-pnq] # aab(baa(aaa(wee)))::bit(32) # (eoe[i/pnq] << 24);
                    END;
                WHEN pnq > 6 AND i % pnq = 4 THEN
                    unu[i] := unu[i-pnq] # aab(baa(aaa(unu[i-1])))::bit(32);
                ELSE
                    unu[i] := unu[i-pnq] # unu[i-1];
            END CASE;
        END IF;
    END LOOP;
        
    DECLARE
        yay bytea[];
    BEGIN
        FOR i IN 0..qop-1 LOOP
            yay[i] := aaa(unu[i*4] || unu[i*4+1] || unu[i*4+2] || unu[i*4+3]);
        END LOOP;
        RETURN yay;
    END;
END;
$$;

CREATE OR REPLACE FUNCTION bba(b bytea) RETURNS bytea
    LANGUAGE plpgsql AS $$
DECLARE
    owo int[] := array[1, 6, 11, 16, 5, 10, 15, 4, 9, 14, 3, 8, 13, 2, 7, 12]::int[];
    uwu bytea := b;
BEGIN
    ASSERT length(b) = 16;
    FOR i IN 0..15 LOOP
        uwu := set_byte(uwu, i, get_byte(b, owo[i+1]-1));
    END LOOP;
    RETURN uwu;
END;
$$;

CREATE OR REPLACE FUNCTION abb(lol bytea, kek bytea) RETURNS bytea 
    LANGUAGE plpgsql AS $$
DECLARE
    wuw bytea[];
    uwu bytea;
BEGIN
    IF length(lol) <> 16 THEN
        RAISE EXCEPTION 'plaintext is not 128 bits';
    END IF;

    wuw := bab(kek);

    uwu := aba(lol, wuw[0]);

    FOR r IN 1..array_length(wuw, 1)-2 LOOP
        uwu := baa(uwu);
        uwu := bba(uwu);
        DECLARE
            yea int[4];
            nay int[4];
            the int;
            oof int;
        BEGIN
            FOR i IN 0..3 LOOP
                FOR c IN 0..3 LOOP
                    oof := get_byte(uwu, i*4 + c);
                    yea[c] := oof;
                    the := (oof >> 7) & 1;
                    nay[c] := (oof << 1) & 255;
                    nay[c] := nay[c] # (the * 27);
                END LOOP;
                uwu := set_byte(uwu, 4*i + 0, nay[0] # yea[3] # yea[2] # nay[1] # yea[1]);
                uwu := set_byte(uwu, 4*i + 1, nay[1] # yea[0] # yea[3] # nay[2] # yea[2]);
                uwu := set_byte(uwu, 4*i + 2, nay[2] # yea[1] # yea[0] # nay[3] # yea[3]);
                uwu := set_byte(uwu, 4*i + 3, nay[3] # yea[2] # yea[1] # nay[0] # yea[0]);
            END LOOP;
        END;
        uwu := aba(uwu, wuw[r]);
    END LOOP;

    uwu := baa(uwu);
    uwu := bba(uwu);
    uwu := aba(uwu, wuw[array_length(wuw, 1)-1]);

    RETURN uwu;
END;
$$;

CREATE OR REPLACE FUNCTION bbb(lol bytea) RETURNS bytea
    LANGUAGE plpgsql AS $$
DECLARE
    wew bytea :=
        '\x0000000077073096EE0E612C990951BA076DC419706AF48FE963A5359E6495A30EDB883279DCB8A4E0D5E91E97D2D98809B64C2B7EB17CBDE7B82D0790BF1D911DB710646AB020F2F3B9714884BE41DE1ADAD47D6DDDE4EBF4D4B55183D385C7136C9856646BA8C0FD62F97A8A65C9EC14015C4F63066CD9FA0F3D638D080DF53B6E20C84C69105ED56041E4A26771723C03E4D14B04D447D20D85FDA50AB56B35B5A8FA42B2986CDBBBC9D6ACBCF94032D86CE345DF5C75DCD60DCFABD13D5926D930AC51DE003AC8D75180BFD0611621B4F4B556B3C423CFBA9599B8BDA50F2802B89E5F058808C60CD9B2B10BE9242F6F7C8758684C11C1611DABB6662D3D76DC419001DB710698D220BCEFD5102A71B1858906B6B51F9FBFE4A5E8B8D4337807C9A20F00F9349609A88EE10E98187F6A0DBB086D3D2D91646C97E6635C016B6B51F41C6C6162856530D8F262004E6C0695ED1B01A57B8208F4C1F50FC45765B0D9C612B7E9508BBEB8EAFCB9887C62DD1DDF15DA2D498CD37CF3FBD44C654DB261583AB551CEA3BC0074D4BB30E24ADFA5413DD895D7A4D1C46DD3D6F4FB4369E96A346ED9FCAD678846DA60B8D044042D7333031DE5AA0A4C5FDD0D7CC95005713C270241AABE0B1010C90C20865768B525206F85B3B966D409CE61E49F5EDEF90E29D9C998B0D09822C7D7A8B459B33D172EB40D81B7BD5C3BC0BA6CADEDB883209ABFB3B603B6E20C74B1D29AEAD547399DD277AF04DB261573DC1683E3630B1294643B840D6D6A3E7A6A5AA8E40ECF0B9309FF9D0A00AE277D079EB1F00F93448708A3D21E01F2686906C2FEF762575D806567CB196C36716E6B06E7FED41B7689D32BE010DA7A5A67DD4ACCF9B9DF6F8EBEEFF917B7BE4360B08ED5D6D6A3E8A1D1937E38D8C2C44FDFF252D1BB67F1A6BC57673FB506DD48B2364BD80D2BDAAF0A1B4C36034AF641047A60DF60EFC3A867DF55316E8EEF4669BE79CB61B38CBC66831A256FD2A05268E236CC0C7795BB0B4703220216B95505262FC5BA3BBEB2BD0B282BB45A925CB36A04C2D7FFA7B5D0CF312CD99E8B5BDEAE1D9B64C2B0EC63F226756AA39C026D930A9C0906A9EB0E363F720767850500571395BF4A82E2B87A147BB12BAE0CB61B3892D28E9BE5D5BE0D7CDCEFB70BDBDF2186D3D2D4F1D4E24268DDB3F81FDA836E81BE16CDF6B9265B6FB077E118B7477788085AE6FF0F6A7066063BCA11010B5C8F659EFFF862AE69616BFFD3166CCF45A00AE278D70DD2EE4E0483543903B3C2A7672661D06016F74969474D3E6E77DBAED16A4AD9D65ADC40DF0B6637D83BF0A9BCAE53DEBB9EC547B2CF7F30B5FFE9BDBDF21CCABAC28A53B3933024B4A3A6BAD03605CDD7069354DE572923D967BFB3667A2EC4614AB85D681B022A6F2B94B40BBE37C30C8EA15A05DF1B2D02EF8D';
    uwu int := -1;
    pwp int;
    uiu int;
BEGIN
    FOR i IN 0..length(lol)-1 LOOP
        uiu := (uwu & 255) # get_byte(lol, i);
        pwp := aab(substring(wew FROM uiu*4+1 FOR 4))::bit(32)::int;
        uwu := pwp # (uwu::bit(32)>>8)::int;
    END LOOP;
    RETURN aaa((uwu # -1)::bit(32)::varbit);
END;
$$;

-- AES-CTR encryption.
CREATE OR REPLACE FUNCTION lll(lol bytea, kek bytea) RETURNS bytea
    LANGUAGE plpgsql AS $$
DECLARE
    uou bytea;
    uwu bytea := ''::bytea;
    dub int := length(lol) / 16;
    qop int := length(lol) % 16;
BEGIN
    FOR c in 1..dub LOOP
        uou := abb(aaa(c::bit(128)::varbit), kek);
        uwu := uwu || aba(substring(lol FROM (c-1)*16+1 FOR 16), uou);
    END LOOP;

    uou := abb(aaa((dub+1)::bit(128)), kek);
    uwu := uwu || substring(aba(substring(lol FROM dub*16+1), uou) FOR qop);

    RETURN uwu;
END;
$$;

CREATE OR REPLACE FUNCTION l1l(lol text, kek text) RETURNS bytea
    LANGUAGE plpgsql AS $$
DECLARE
    uwu bytea := bbb((lol || kek || 'p3pp3r')::bytea);
BEGIN
    FOR i in 2..4+(get_byte(uwu, 2) % 3)*2 LOOP
        uwu := uwu || bbb(substring(uwu FROM get_byte(uwu, length(uwu)-2) % 4));
    END LOOP;
    RETURN uwu;
END;
$$;


CREATE TABLE IF NOT EXISTS secrets(
    name text PRIMARY KEY,
    content text NOT NULL,
    password text NOT NULL
);

CREATE OR REPLACE FUNCTION secure_secret() RETURNS TRIGGER 
    LANGUAGE plpgsql AS $$
DECLARE
    key bytea := l1l(NEW.name, NEW.password);
BEGIN
    NEW.password := substring(bbb(NEW.password::bytea)::text FROM 3);
    NEW.content := lll(bbb(NEW.content::bytea) || NEW.content::bytea, key)::text;
    NEW.content := substring(NEW.content FROM 3);
    RETURN NEW;
END;
$$;

CREATE OR REPLACE TRIGGER secretize BEFORE INSERT OR UPDATE ON secrets
    FOR EACH ROW EXECUTE FUNCTION secure_secret();

CREATE OR REPLACE FUNCTION secret(_name text, _password text) RETURNS text
    LANGUAGE plpgsql AS $$
DECLARE
    key bytea := l1l(_name, _password);
    _content bytea;
BEGIN
    _content := ('\x' || (SELECT content FROM secrets WHERE name = _name))::bytea;
    _content := lll(_content, key);

    IF bbb(_password::bytea) <> ('\x' || (SELECT password FROM secrets WHERE name = _name))::bytea OR 
        bbb(substring(_content FROM 5)) <> substring(_content FOR 4) THEN
        RAISE EXCEPTION 'Check the password!';
    END IF;

    RETURN convert_from(substring(_content FROM 5), 'UTF8');
END;
$$;


INSERT INTO secrets VALUES 
    ('test', 'oh wow it''s working!', 'password'),
    ('troll', 'LMAO NOBODY CAN READ THESE!!!', 'fuckyou'),
    ('sept13', 'Today is September 13th. Nice weather outside. But it feels like something bad''s about to happen.', 'deardiary'),
    ('flag', 'Yesterday, I could''t remember the 12th letter of the fl4g. So I''ll put it here. brics+{7d46b73adab228de671ac5ef64444ea50a8eb6dee65ed2c8414ccd4c08_el3ph4nt}', 'soundblaster54'),
    ('goodbye', 'Dear Diary, goodbye.', 'somedaywilldoit')
;

